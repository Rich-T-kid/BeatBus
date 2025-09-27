import { clsx } from 'clsx'

const LoadingSpinner = ({ size = 'md', className }) => {
  const sizes = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8', 
    lg: 'w-12 h-12',
    xl: 'w-16 h-16'
  }
  
  return (
    <div className={clsx('spinner', sizes[size], className)} />
  )
}

export default LoadingSpinner