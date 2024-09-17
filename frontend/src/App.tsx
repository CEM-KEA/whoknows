import { BrowserRouter, Route, Routes } from "react-router-dom";
import Nav from "./components/Nav";
import Search from "./views/Search";
import Login from "./views/Login";
import { useEffect, useState } from "react";
import {
  getUserCookieAccept,
  removeJWTTokenFromCookies,
  setJWTTokenInCookies,
  setUserCookieAccept
} from "./helpers/cookieHelpers";
import CookieBanner from "./components/CookieBanner";

function App() {
  const [loggedIn, setLoggedIn] = useState<boolean>(false);
  const [showCookieBanner, setShowCookieBanner] = useState<boolean>(true);

  useEffect(() => {
    if (getUserCookieAccept()) {
      setShowCookieBanner(false);
    }
  }, []);

  function logOut() {
    setLoggedIn(false);
    removeJWTTokenFromCookies();
  }

  function logIn(jwt_token: string) {
    setJWTTokenInCookies(jwt_token);
    setLoggedIn(true);
  }

  return (
    <>
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
        {showCookieBanner && (
          <CookieBanner
            onChoice={(choice) => {
              setShowCookieBanner(false);
              setUserCookieAccept(choice);
            }}
          />
        )}
      </BrowserRouter>
    </>
  );
}

export default App;
