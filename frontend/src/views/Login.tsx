import { FormEventHandler, useState } from "react";
import PageLayout from "../components/PageLayout";
import { ILoginRequest, ILoginResponse } from "../types/auth.types";
import { apiPost } from "../utils/apiUtils";
import { useNavigate } from "react-router-dom";
import LoadingSpinner from "../components/LoadingSpinner";
import toast from "react-hot-toast";

interface LoginProps {
  onLogIn: (token: string) => void;
}

function Login(props: Readonly<LoginProps>) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleLogin: FormEventHandler = (e) => {
    e.preventDefault();
    setLoading(true);
    apiPost<ILoginRequest, ILoginResponse>("/login", { username: username, password })
      .then((data) => {
        setLoading(false);
        if (data.require_password_change) {
          toast.error("You need to change your password.");
          navigate("/change-password");
          return;
        }
        props.onLogIn(data.token);
        toast.success("Logged in successfully.");
        navigate("/");
      })
      .catch((error) => {
        toast.error(error.message);
        setLoading(false);
      });
  };

  return (
    <PageLayout>
      <div className="flex justify-center items-center h-[calc(100vh-180px)]">
        {loading && <LoadingSpinner size={100} />}
        {!loading && (
          <form
            className="flex flex-col gap-4 border-2 px-20 py-8 rounded-xl w-1/2"
            onSubmit={handleLogin}
          >
            <h2 className="text-2xl font-semibold">Log in</h2>
            <input
              id="login-username"
              type="text"
              placeholder="Username"
              className="border-2 p-2 w-full rounded outline-none caret-blue-500 text-xl"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            <input
              id="login-password"
              type="password"
              placeholder="Password"
              className="border-2 p-2 w-full rounded outline-none caret-blue-500 text-xl mt-2"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            <div className="flex justify-center">
              <button
                id="login-button"
                className="border-2 rounded w-1/2 mt-2 p-2 bg-blue-500 text-white font-semibold hover:brightness-90 text-xl"
              >
                Log in
              </button>
            </div>
          </form>
        )}
      </div>
    </PageLayout>
  );
}

export default Login;
