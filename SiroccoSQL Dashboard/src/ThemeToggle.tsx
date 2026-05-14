import { Moon, Sun } from "lucide-react";
import { useTheme } from "./hooks/useTheme";

export default function ThemeToggle() {
    const { dark, setDark } = useTheme();

    return (
        <button
            onClick={() => setDark((prev) => !prev)}
            className="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-800 transition"
        >
            {dark ? <Sun size={18} /> : <Moon size={18} />}
        </button>
    );
}