import { BrowserRouter, Route, Routes } from "react-router-dom";
import Nav from "./components/Nav";
import Search from "./views/Search";
import Login from "./views/Login";
import { useState } from "react";
import { removeJWTTokenFromCookies, setJWTTokenInCookies } from "./helpers/cookieHelpers";
import Weather from "./views/Weather";

function App() {
  const [loggedIn, setLoggedIn] = useState<boolean>(false);

  function logOut() {
    setLoggedIn(false);
    removeJWTTokenFromCookies();
  }

  function logIn(jwt_token: string) {
    setJWTTokenInCookies(jwt_token);
    setLoggedIn(true);
  }

  return (
    <BrowserRouter>
      <Nav
        loggedIn={loggedIn}
        onLogOut={logOut}
      />
      <Routes>
        <Route
          path="/"
          element={<Search />}
        />
        <Route
          path="/weather"
          element={<Weather />}
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
  );
}

export default App;
