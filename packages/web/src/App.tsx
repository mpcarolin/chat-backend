import './App.css'
import { ChatDemoContainer } from './components/containers/ChatDemoContainer/ChatDemoContainer'
import { Header } from './components/common/Header'

function App() {

  return (
    <div className="app-root">
      <Header />
      <ChatDemoContainer />
    </div>
  )
}

export default App
