import {
    useEffect,
    useRef,
    useState,
} from 'react';

export function useCanvasViewport() {
    const viewportRef =
        useRef<HTMLDivElement>(null);

    const [scale, setScale] = useState(2);
    const [pan, setPan] = useState({
        x: 0,
        y: 0,
    });

    const dragging = useRef(false);
    const didDrag = useRef(false);

    const lastMousePosition = useRef({
        x: 0,
        y: 0,
    });

    const scaleRef = useRef(2);
    const panRef = useRef({
        x: 0,
        y: 0,
    });

    const targetPanRef = useRef({
        x: 0,
        y: 0,
    });

    useEffect(() => {
        const handleMouseUp = () => {
            dragging.current = false;
        };

        window.addEventListener(
            'mouseup',
            handleMouseUp
        );

        return () => {
            window.removeEventListener(
                'mouseup',
                handleMouseUp
            );
        };
    }, []);

    const handleWheel = (
        event: React.WheelEvent<HTMLDivElement>
    ) => {
        event.preventDefault();

        const viewport = viewportRef.current;

        if (!viewport) return;

        const rect = viewport.getBoundingClientRect();

        const mouseX = event.clientX - rect.left;
        const mouseY = event.clientY - rect.top;

        const centerX = rect.width / 2;
        const centerY = rect.height / 2;

        const currentScale = scaleRef.current;
        const currentPan = panRef.current;

        const newScale = Math.min(
            50,
            Math.max(
                1.5,
                currentScale *
                (event.deltaY < 0 ? 1.2 : 1 / 1.2)
            )
        );

        if (newScale === currentScale) {
            return;
        }

        const worldX =
            (mouseX - centerX - currentPan.x) /
            currentScale;

        const worldY =
            (mouseY - centerY - currentPan.y) /
            currentScale;

        const newPan = {
            x: mouseX - centerX - worldX * newScale,
            y: mouseY - centerY - worldY * newScale,
        };

        scaleRef.current = newScale;
        panRef.current = newPan;

        setScale(newScale);
        setPan(newPan);
    };

    const handleMouseDown = (
        event: React.MouseEvent<HTMLDivElement>
    ) => {
        if (event.button !== 0) return;
        if (event.shiftKey) return;

        event.preventDefault();

        dragging.current = true;
        didDrag.current = false;

        lastMousePosition.current = {
            x: event.clientX,
            y: event.clientY,
        };
    };

    const handleMouseMove = (
        event: React.MouseEvent<HTMLDivElement>
    ) => {

        if (!dragging.current) return;

        const dx =
            event.clientX -
            lastMousePosition.current.x;

        const dy =
            event.clientY -
            lastMousePosition.current.y;

        if (Math.abs(dx) > 1 || Math.abs(dy) > 1) {
            didDrag.current = true;
        }

        const nextPan = {
            x:
                panRef.current.x +
                dx,
            y:
                panRef.current.y +
                dy,
        };

        panRef.current = nextPan;
        targetPanRef.current =
            nextPan;

        setPan(nextPan);

        lastMousePosition.current = {
            x: event.clientX,
            y: event.clientY,
        };
    };

    return {
        viewportRef,
        scale,
        pan,
        didDrag,
        handleWheel,
        handleMouseDown,
        handleMouseMove,
    };
}