import { BrowserRouter, Route, Routes } from "react-router-dom";
import Nav from "./components/Nav";
import Search from "./views/Search";
import Login from "./views/Login";
import { useEffect, useState } from "react";
import {
  getJWTTokenFromCookies,
  removeJWTTokenFromCookies,
  setJWTTokenInCookies
} from "./helpers/cookieHelpers";
import Weather from "./views/Weather";
import Register from "./views/Register";
import toast, { Toaster } from "react-hot-toast";
import { apiGetVoid } from "./utils/apiUtils";
import ChangePassword from "./views/ChangePassword";
import DataLeakBanner from "./components/DataLeakBanner";

function App() {
  const [loggedIn, setLoggedIn] = useState<boolean>(false);
  const [showBanner, setShowBanner] = useState<boolean>(true);

  // right now, as the jwt token is not really used, we just check if it exists to see if the user is logged in
  // ideally, we would also check if the token is still valid with the backend
  useEffect(() => {
    const jwt_token = getJWTTokenFromCookies();
    if (!jwt_token || jwt_token === "") {
      setLoggedIn(false);
    } else {
      setLoggedIn(true);
    }
  }, []);

  function logOut() {
    setLoggedIn(false);
    void apiGetVoid("/logout", true)
      .catch((e) => toast.error(e.message))
      .finally(() => removeJWTTokenFromCookies());
  }

  function logIn(jwt_token: string) {
    setJWTTokenInCookies(jwt_token);
    setLoggedIn(true);
  }

  return (
    <BrowserRouter>
      {showBanner && <DataLeakBanner onClose={() => setShowBanner(false)} />}
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
        {!loggedIn && (
          <>
            <Route
              path="/register"
              element={<Register logIn={logIn} />}
            />
            <Route
              path="/login"
              element={<Login onLogIn={logIn} />}
            />
          </>
        )}
        <Route
          path="/change-password"
          element={<ChangePassword />}
        />
      </Routes>
      <Toaster position="bottom-left" />
    </BrowserRouter>
  );
}

export default App;
