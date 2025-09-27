import { clsx } from 'clsx'

const Card = ({ children, className, ...props }) => {
  return (
    <div
      className={clsx(
        'rounded-xl border border-gray-200 bg-white p-6 shadow-sm',
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
}

export default Card