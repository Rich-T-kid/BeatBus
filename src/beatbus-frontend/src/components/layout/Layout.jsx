import { Outlet } from 'react-router-dom'
import Header from './Header'

const Layout = () => {
  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main>
        <Outlet />
      </main>
    </div>
  )
}

export default Layout