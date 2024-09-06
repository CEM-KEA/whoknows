import { BrowserRouter, Route, Routes } from "react-router-dom"
import Nav from "./components/Nav"

function App() {
  return (
    <>
      <BrowserRouter>
        <Nav loggedIn={false} />
        <Routes>
          <Route path="/" element={<div />} /> 
          <Route path="/weather" element={<div />} />
          <Route path="/register" element={<div />} />
          <Route path="/login" element={<div />} />
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
