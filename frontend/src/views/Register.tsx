import { FormEventHandler, useState } from "react";
import type { IRegisterRequest } from "../types/auth.types";
import PageLayout from "../components/PageLayout";
import { apiPost } from "../utils/apiUtils";
import { useNavigate } from "react-router-dom";
import LoadingSpinner from "../components/LoadingSpinner";

function Register() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit: FormEventHandler = async (e) => {
    e.preventDefault();
    const registerData: IRegisterRequest = {
      username,
      email,
      password
    };
    setLoading(true);
    apiPost<IRegisterRequest, void>("/register", registerData)
      .then(() => {
        setLoading(false);
        navigate("/login");
      })
      .catch((error) => {
        console.error(error);
        setLoading(false);
        // Maybe do a toast here
      });
  };

  return (
    <PageLayout>
      <div className="flex items-center justify-center mt-20 p-8">
        {loading && (
          <div className="mt-20">
            <LoadingSpinner size={100} />
          </div>
        )}
        {!loading && (
          <form
            onSubmit={handleSubmit}
            className="flex flex-col gap-12 lg:w-1/2 border-2 px-20 py-8 rounded bg-blue-200 bg-opacity-40"
          >
            <h1 className="text-2xl font-semibold">Register new user</h1>
            <input
              className="border-2 p-2 w-full rounded outline-none caret-blue-500 text-xl"
              type="text"
              placeholder="Username"
              value={username}
              min={3}
              max={100}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            <input
              className="border-2 p-2 w-full rounded outline-none caret-blue-500 text-xl"
              type="email"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <input
              className="border-2 p-2 w-full rounded outline-none caret-blue-500 text-xl"
              type="password"
              placeholder="Password"
              value={password}
              minLength={6}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            <div className="flex gap-1 w-full justify-between mt-8">
              <button
                className="border rounded bg-blue-50 text-blue-500 hover:bg-blue-100 font-semibold p-2"
                onClick={() => {
                  setUsername("");
                  setEmail("");
                  setPassword("");
                }}
                type="button"
              >
                Clear
              </button>
              <button className="border rounded bg-blue-500 hover:brightness-90 text-white font-semibold p-2">
                Register
              </button>
            </div>
          </form>
        )}
      </div>
    </PageLayout>
  );
}

export default Register;
