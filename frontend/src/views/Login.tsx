import { useState } from "react";
import PageLayout from "../components/PageLayout";
import { ILoginRequest, ILoginResponse } from "../types/types";
import { apiPost } from "../utils/apiUtils";
import { useNavigate } from "react-router-dom";

interface LoginProps {
  onLogIn: (token: string) => void;
}

function Login(props: Readonly<LoginProps>) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  function handleLogin() {
    apiPost<ILoginRequest, ILoginResponse>("/login", { username, password })
      .then((data) => {
        props.onLogIn(data.token);
        navigate("/");
      })
      .catch((error) => {
        console.error(error);
        // Maybe do a toast here
      });
  }

  return (
    <PageLayout>
      <div className="flex justify-center items-center h-[calc(100vh-180px)]">
        <div className="flex flex-col gap-4 border-2 px-20 py-8 rounded-xl w-1/2">
          <h2 className="text-2xl font-semibold">Log in</h2>
          <input
            id="login-username"
            type="text"
            placeholder="Username"
            className="border-2 p-2 w-full rounded outline-none caret-blue-500 text-xl"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            id="login-password"
            type="password"
            placeholder="Password"
            className="border-2 p-2 w-full rounded outline-none caret-blue-500 text-xl mt-2"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <div className="flex justify-center">
            <button
              id="login-button"
              onClick={handleLogin}
              className="border-2 rounded w-1/2 mt-2 p-2 bg-blue-500 text-white font-semibold hover:brightness-90 text-xl"
            >
              Log in
            </button>
          </div>
        </div>
      </div>
    </PageLayout>
  );
}

export default Login;
