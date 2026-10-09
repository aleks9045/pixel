import {
    useCallback,
    useEffect,
    useRef,
} from 'react';

import config from '../../../config';
import { Ws } from '../../../proto/pixel';

type PixelRegion = {
    x1: number;
    x2: number;
    y1: number;
    y2: number;
    color: number;
};

export function usePixelWebSocket(
    onPixel: (
        x: number,
        y: number,
        color: number
    ) => void,
    onPixelRegion: (region: PixelRegion) => void,
    onError: (message: string) => void
) {
    const wsRef =
        useRef<WebSocket | null>(null);

    useEffect(() => {
        const ws =
            new WebSocket(
                `${config.host}/ws/connect`
            );

        ws.binaryType =
            'arraybuffer';

        ws.onopen = () => {
            console.log(
                'WebSocket connected'
            );
        };

        ws.onmessage = event => {
            const message =
                Ws.decode(
                    new Uint8Array(
                        event.data
                    )
                );

            if (message.pixel) {
                console.log(message.pixel)
                onPixel(
                    message.pixel.x,
                    message.pixel.y,
                    message.pixel.color
                );

                return;
            }
            if (message.pixelRegion) {
                onPixelRegion({
                    x1: message.pixelRegion.x1,
                    x2: message.pixelRegion.x2,
                    y1: message.pixelRegion.y1,
                    y2: message.pixelRegion.y2,
                    color: message.pixelRegion.color,
                });
                return;
            }

            if (message.pixelError) {
                onError(
                    message.pixelError.message
                );
            }
            if (message.pixelRegionError) {
                onError(
                    message.pixelRegionError.message
                );
            }
        };

        ws.onclose = () => {
            console.log(
                'WebSocket disconnected'
            );

            if (
                wsRef.current === ws
            ) {
                wsRef.current = null;
            }
        };

        ws.onerror = error => {
            console.error(
                'WebSocket error',
                error
            );
        };

        wsRef.current = ws;

        return () => {
            if (
                wsRef.current === ws
            ) {
                wsRef.current = null;
            }

            ws.close();
        };
    }, [onPixel, onError]);

    const sendPixel =
        useCallback(
            (
                x: number,
                y: number,
                color: number
            ) => {
                const ws =
                    wsRef.current;

                if (
                    !ws ||
                    ws.readyState !==
                    WebSocket.OPEN
                ) {
                    onError(
                        'WebSocket connection is not available'
                    );

                    return;
                }

                const message =
                    Ws.encode({
                        pixel: {
                            x,
                            y,
                            color,
                        },
                    }).finish();

                ws.send(message);
            },
            [onError]
        );

    const sendPixelRegion =
        useCallback(
            (
                x1: number,
                x2: number,
                y1: number,
                y2: number,
                color: number
            ) => {
                const ws =
                    wsRef.current;

                if (
                    !ws ||
                    ws.readyState !==
                    WebSocket.OPEN
                ) {
                    onError(
                        'WebSocket connection is not available'
                    );

                    return;
                }

                const message = Ws.encode({
                    pixelRegion: {
                        x1,
                        x2,
                        y1,
                        y2,
                        color,
                    },
                }).finish();
                +

                    ws.send(message);
            },
            [onError]
        );
    return {
        sendPixel,
        sendPixelRegion,
    };
}