import React from 'react';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    isLoading?: boolean;
    variant?: 'primary' | 'secondary' | 'outline';
}

export const Button: React.FC<ButtonProps> = ({ 
    children, 
    isLoading, 
    variant = 'primary', 
    className = '', 
    disabled,
    ...props 
}) => {
    const baseStyles = "w-full py-2 px-4 rounded-lg font-semibold transition-all duration-200 flex justify-center items-center";
    
    const variants = {
        primary: "bg-indigo-600 text-white hover:bg-indigo-700 shadow-md hover:shadow-lg",
        secondary: "bg-gray-200 text-gray-800 hover:bg-gray-300",
        outline: "border-2 border-indigo-600 text-indigo-600 hover:bg-indigo-50"
    };

    return (
        <button
            disabled={isLoading || disabled}
            className={`${baseStyles} ${variants[variant]} ${isLoading || disabled ? 'opacity-70 cursor-not-allowed' : ''} ${className}`}
            {...props}
        >
            {isLoading ? (
                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
            ) : null}
            {children}
        </button>
    );
};