import { BrowserRouter, Route, Routes } from "react-router-dom";
import Nav from "./components/Nav";
import Search from "./views/Search";
import Login from "./views/Login";
import Cookies from "universal-cookie";
import { ILoginSession } from "./types/types";
import { useState } from "react";

function App() {
  const cookies = new Cookies();
  const [loginSession, setLoginSession] = useState<ILoginSession | null>(
    cookies.get("jwt_authorization")
  );

  function logOut() {
    setLoginSession(null);
    cookies.remove("jwt_authorization");
  }

  function logIn(jwt_token: string, username: string) {
    cookies.set("jwt_authorization", jwt_token);
    setLoginSession({ username });
  }

  return (
    <>
      <BrowserRouter>
        <Nav
          loggedIn={!!loginSession}
          onLogOut={logOut}
        />
        <Routes>
          <Route
            path="/"
            element={<Search />}
          />
          <Route
            path="/weather"
            element={<div />}
          />
          <Route
            path="/register"
            element={<div />}
          />
          <Route
            path="/login"
            element={<Login onLogIn={logIn} />}
          />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;
