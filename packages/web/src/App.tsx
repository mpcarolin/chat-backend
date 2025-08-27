import './App.css'
import { ChatDemoContainer } from './components/containers/ChatDemoContainer/ChatDemoContainer'

function App() {

  return (
    <div className="app-root">
      <Header />
      <ChatDemoContainer />
    </div>
  )
}

const Header = () => (
  <div className="app-header">
    <img
      width={132}
      height={32}
      src="https://www.neudesic.com/wp-content/uploads/neudesic-wht-logo-x2.png"
      alt="Neudesic Logo"
    />
  </div>
)

export default App
