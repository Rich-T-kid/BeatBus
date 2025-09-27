import { forwardRef } from 'react'
import { clsx } from 'clsx'

const Input = forwardRef(({ 
  label, 
  error, 
  className, 
  ...props 
}, ref) => {
  return (
    <div className="space-y-2">
      {label && (
        <label className="label">
          {label}
        </label>
      )}
      <input
        ref={ref}
        className={clsx(
          'input',
          error && 'border-red-500 focus:border-red-500 focus:ring-red-500',
          className
        )}
        {...props}
      />
      {error && (
        <p className="text-sm text-red-600">{error}</p>
      )}
    </div>
  )
})

Input.displayName = 'Input'

export default Input