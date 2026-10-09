import { useEffect, useState } from "react";
import { getPixelField } from "../../../api/pixel";


export function usePixelField() {
    const [pixelColors, setPixelColors] =
        useState<Uint8Array>();

    useEffect(() => {
        getPixelField()
            .then(field => {
                setPixelColors(field.pixelColors);
            })
            .catch(console.error);
    }, []);

    return {
        pixelColors,
    };
}