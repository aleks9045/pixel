import {
    useCallback,
    useRef,
} from 'react';
import config from '../../../config';
import { colors } from '../constants/colors';

export function usePixelCanvas() {
    const canvasRef =
        useRef<HTMLCanvasElement>(null);

    const drawCanvas = useCallback(
        (pixelColors: Uint8Array) => {
            const canvas =
                canvasRef.current;

            if (!canvas) return;

            canvas.width =
                config.fieldWidth;

            canvas.height =
                config.fieldHeight;

            const context =
                canvas.getContext('2d');

            if (!context) return;

            const image =
                context.createImageData(
                    config.fieldWidth,
                    config.fieldHeight
                );

            for (
                let i = 0;
                i < pixelColors.length;
                i++
            ) {
                const color =
                    colors[pixelColors[i]];

                if (!color) continue;

                const hex =
                    color.slice(1);

                image.data[i * 4] =
                    parseInt(
                        hex.slice(0, 2),
                        16
                    );

                image.data[i * 4 + 1] =
                    parseInt(
                        hex.slice(2, 4),
                        16
                    );

                image.data[i * 4 + 2] =
                    parseInt(
                        hex.slice(4, 6),
                        16
                    );

                image.data[i * 4 + 3] =
                    255;
            }

            context.putImageData(
                image,
                0,
                0
            );
        },
        []
    );

    const drawPixel = useCallback(
        (
            x: number,
            y: number,
            color: number
        ) => {
            const canvas =
                canvasRef.current;

            if (!canvas) return;

            const pixelColor =
                colors[color];

            if (!pixelColor) return;

            const context =
                canvas.getContext('2d');

            if (!context) return;

            context.fillStyle =
                pixelColor;

            context.fillRect(
                x,
                y,
                1,
                1
            );
        },
        []
    );

    const getPixelCoordinates =
        useCallback(
            (
                event: React.MouseEvent<HTMLCanvasElement>
            ) => {
                const canvas =
                    canvasRef.current;

                if (!canvas) return null;

                const rect =
                    canvas.getBoundingClientRect();

                const x = Math.floor(
                    ((event.clientX -
                        rect.left) /
                        rect.width) *
                    config.fieldWidth
                );

                const y = Math.floor(
                    ((event.clientY -
                        rect.top) /
                        rect.height) *
                    config.fieldHeight
                );

                if (
                    x < 0 ||
                    x >= config.fieldWidth ||
                    y < 0 ||
                    y >= config.fieldHeight
                ) {
                    return null;
                }

                return {
                    x,
                    y,
                };
            },
            []
        );


    return {
        canvasRef,
        drawCanvas,
        drawPixel,
        getPixelCoordinates
    };
}