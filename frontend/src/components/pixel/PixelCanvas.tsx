import {
    useCallback,
    useEffect,
    useRef,
    useState,
} from 'react';

import config from '../../config';
import { ColorPalette } from './ColorPalette';
import { useCanvasViewport } from './use/useCanvasViewport';
import { usePixelCanvas } from './use/usePixelCanvas';
import { usePixelField } from './use/usePixelField';
import { usePixelWebSocket } from './use/usePixelWebSocket';
import { ErrorToast } from '../ErrorToast';
import { PixelTimeout } from './PixelTimeout';
import { PixelCords } from './PixelCords';
import { useAuth } from '../../context/AuthContext';
import { PixelAuthor } from './PixelAuthor';
import { usePathname } from '../../hooks/usePathname';

type PixelRegion = {
    x1: number;
    x2: number;
    y1: number;
    y2: number;
};

export function PixelCanvas() {

    const { user, loading } = useAuth()


    const [pixelPosition, setPixelPosition] =
        useState<{ x: number; y: number } | null>(null);

    const [pixelPlaced, setPixelPlaced] =
        useState(0);

    const [canPlacePixel, setCanPlacePixel] =
        useState(true);

    const [selectedColor, setSelectedColor] =
        useState(2);

    const [error, setError] =
        useState<string | null>(null);

    const [showAuthor, setShowAuthor] =
        useState(false);

    const [authorCardPosition, setAuthorCardPosition] =
        useState<{ x: number; y: number } | null>(null);

    const [authorCoordinates, setAuthorCoordinates] =
        useState<{ x: number; y: number } | null>(null);

    const [pixelRegion, setPixelRegion] =
        useState<PixelRegion | null>(null);

    const selectingRegion = useRef(false);

    const pixelRegionRef =
        useRef<PixelRegion | null>(null);

    const suppressNextClick = useRef(false);

    const {
        canvasRef,
        drawCanvas,
        drawPixel,
        getPixelCoordinates,
    } = usePixelCanvas();

    const {
        pixelColors,
    } = usePixelField();


    const handlePixelRegion = useCallback(
        ({ x1, x2, y1, y2, color }: {
            x1: number;
            x2: number;
            y1: number;
            y2: number;
            color: number;
        }) => {
            const minX = Math.min(x1, x2);
            const maxX = Math.max(x1, x2);
            const minY = Math.min(y1, y2);
            const maxY = Math.max(y1, y2);

            for (let y = minY; y <= maxY; y++) {
                for (let x = minX; x <= maxX; x++) {
                    drawPixel(x, y, color);
                }
            }
        },
        [drawPixel]
    );

    const {
        sendPixel,
        sendPixelRegion
    } = usePixelWebSocket(
        drawPixel,
        handlePixelRegion,
        setError
    );

    const {
        viewportRef,
        scale,
        pan,
        didDrag,
        handleWheel,
        handleMouseDown,
        handleMouseMove,
    } = useCanvasViewport();

    useEffect(() => {
        if (!pixelColors) return;

        drawCanvas(pixelColors);
    }, [
        pixelColors,
        drawCanvas,
    ]);

    const handleContextMenu = (
        event: React.MouseEvent<HTMLCanvasElement>
    ) => {
        event.preventDefault();

        if (!canPlacePixel) return;

        const coordinates =
            getPixelCoordinates(event);

        if (!coordinates) return;

        setAuthorCardPosition({
            x: event.clientX,
            y: event.clientY,
        });

        setAuthorCoordinates(coordinates);
        setShowAuthor(true);
    };


    const handleDrawPixel = (
        event: React.MouseEvent<HTMLCanvasElement>
    ) => {
        event.preventDefault();

        if (suppressNextClick.current) {
            suppressNextClick.current = false;
            return;
        }

        if (didDrag.current) return;
        if (!canPlacePixel) return;
        if (pixelRegion) return;

        const coordinates =
            getPixelCoordinates(event);

        if (!coordinates) return;

        setShowAuthor(false);

        const {
            x,
            y,
        } = coordinates;

        drawPixel(
            x,
            y,
            selectedColor
        );

        sendPixel(
            x,
            y,
            selectedColor
        );

        setPixelPlaced(value => value + 1);
    };

    if (loading || !user) {
        return null;
    }

    const handleRegionMouseDown = (
        event: React.MouseEvent<HTMLCanvasElement>
    ) => {
        if (!user.isAdmin) return;
        if (!event.shiftKey) return;
        if (event.button !== 0) return;

        const coordinates =
            getPixelCoordinates(event);

        if (!coordinates) return;

        event.preventDefault();
        suppressNextClick.current = true;
        selectingRegion.current = true;

        const region = {
            x1: coordinates.x,
            y1: coordinates.y,
            x2: coordinates.x,
            y2: coordinates.y,
        };

        pixelRegionRef.current = region;
        setPixelRegion(region);

        setShowAuthor(false);
    };

    const handleCanvasMouseMove = (
        event: React.MouseEvent<HTMLCanvasElement>
    ) => {
        const coordinates =
            getPixelCoordinates(event);


        setPixelPosition(coordinates);

        if (!user.isAdmin) return;
        if (!selectingRegion.current) return;
        if (!coordinates) return;

        const currentRegion =
            pixelRegionRef.current;
        if (!currentRegion) return;

        const region = {
            x1: currentRegion.x1,
            x2: coordinates.x,
            y1: currentRegion.y1,
            y2: coordinates.y,
        };


        pixelRegionRef.current = region;
        setPixelRegion(region);
    };

    const handleRegionMouseUp = () => {
        if (!user.isAdmin) return;
        if (!selectingRegion.current) return;

        selectingRegion.current = false;

        const region =
            pixelRegionRef.current;

        pixelRegionRef.current = null;
        setPixelRegion(null);

        console.log(
            region
        );
        if (!region) return;

        const x1 = Math.min(
            region.x1,
            region.x2
        );

        const x2 = Math.max(
            region.x1,
            region.x2
        );

        const y1 = Math.min(
            region.y1,
            region.y2
        );

        const y2 = Math.max(
            region.y1,
            region.y2
        );

        sendPixelRegion(
            x1,
            x2,
            y1,
            y2,
            selectedColor
        );
    };

    const { navigate } = usePathname();

    return (

        <div className="pixel-editor">

            {user.isAdmin && (
                <button
                    className="admin-open-button"
                    onClick={() => navigate('/admin')}
                >
                    Админ-панель
                </button>
            )}

            <ColorPalette
                selectedColor={
                    selectedColor
                }
                onSelect={
                    setSelectedColor
                }
            />

            <div
                ref={viewportRef}
                className="pixel-viewport"
                onWheel={handleWheel}
                onMouseDown={handleMouseDown}
                onMouseMove={handleMouseMove}
                onMouseUp={handleRegionMouseUp}
                onMouseLeave={() => setPixelPosition(null)}
            >
                <canvas
                    ref={canvasRef}
                    onClick={handleDrawPixel}
                    onMouseMove={handleCanvasMouseMove}
                    onMouseDown={handleRegionMouseDown}
                    onContextMenu={handleContextMenu}
                    className="pixel-canvas"
                    style={{
                        width: `${config.fieldWidth * scale}px`,
                        height: `${config.fieldHeight * scale}px`,
                        left: `calc(50% - ${config.fieldWidth * scale / 2}px + ${pan.x}px)`,
                        top: `calc(50% - ${config.fieldHeight * scale / 2}px + ${pan.y}px)`,
                    }}
                />
            </div>

            {pixelPlaced && (<PixelTimeout
                time={user.timeLimit}
                trigger={pixelPlaced}
                onChange={setCanPlacePixel}
            />)}

            {error && (
                <ErrorToast
                    message={error}
                    onClose={() => setError(null)}
                />
            )}

            {pixelPosition && (
                <PixelCords
                    x={pixelPosition.x}
                    y={pixelPosition.y}
                />
            )}

            {showAuthor && authorCoordinates && authorCardPosition && (
                <PixelAuthor
                    x={authorCoordinates.x}
                    y={authorCoordinates.y}
                    left={authorCardPosition.x}
                    top={authorCardPosition.y}
                />
            )}

            {user.isAdmin && pixelRegion && (
                <div
                    className="pixel-region"
                    style={{
                        left: `calc(
                50% -
                ${config.fieldWidth * scale / 2}px +
                ${pan.x + Math.min(
                            pixelRegion.x1,
                            pixelRegion.x2
                        ) * scale}px
            )`,
                        top: `calc(
                50% -
                ${config.fieldHeight * scale / 2}px +
                ${pan.y + Math.min(
                            pixelRegion.y1,
                            pixelRegion.y2
                        ) * scale}px
            )`,
                        width: `${(
                            Math.abs(
                                pixelRegion.x2 -
                                pixelRegion.x1
                            ) + 1
                        ) * scale
                            }px`,
                        height: `${(
                            Math.abs(
                                pixelRegion.y2 -
                                pixelRegion.y1
                            ) + 1
                        ) * scale
                            }px`,
                    }}
                />
            )}
        </div>
    );
}