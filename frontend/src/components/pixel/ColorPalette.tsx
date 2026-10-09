import { colors } from './constants/colors';

type Props = {
    selectedColor: number;
    onSelect: (color: number) => void;
};

export function ColorPalette({
    selectedColor,
    onSelect,
}: Props) {
    return (
        <div className="color-menu">
            {colors.map((color, index) => (
                <button
                    key={index}
                    className={`color-button ${selectedColor === index ? 'selected' : ''
                        }`}
                    style={{ backgroundColor: color }}
                    onClick={() => onSelect(index)}
                />
            ))}
        </div>
    );
}