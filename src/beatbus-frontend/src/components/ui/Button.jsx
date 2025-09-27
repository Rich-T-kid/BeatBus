import { forwardRef } from 'react'
import { clsx } from 'clsx'

const Button = forwardRef(({ 
  children, 
  variant = 'primary', 
  size = 'md', 
  disabled, 
  loading, 
  className, 
  ...props 
}, ref) => {
  const baseStyles = 'inline-flex items-center justify-center font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50'
  
  const variants = {
    primary: 'bg-primary-600 text-white hover:bg-primary-700 focus:ring-primary-500',
    secondary: 'bg-secondary-600 text-white hover:bg-secondary-700 focus:ring-secondary-500',
    outline: 'border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-primary-500',
    ghost: 'text-gray-600 hover:bg-gray-100 hover:text-gray-900 focus:ring-gray-500',
    danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500'
  }
  
  const sizes = {
    sm: 'px-3 py-2 text-sm rounded-md',
    md: 'px-4 py-2 text-sm rounded-lg',
    lg: 'px-6 py-3 text-base rounded-lg',
    xl: 'px-8 py-4 text-lg rounded-xl'
  }
  
  return (
    <button
      ref={ref}
      disabled={disabled || loading}
      className={clsx(
        baseStyles,
        variants[variant],
        sizes[size],
        className
      )}
      {...props}
    >
      {loading && (
        <div className="spinner w-4 h-4 mr-2" />
      )}
      {children}
    </button>
  )
})

Button.displayName = 'Button'

export default Button