import { FormEventHandler, useMemo, useState } from "react";
import type { IRegisterRequest } from "../types/auth.types";
import PageLayout from "../components/PageLayout";
import { apiPost } from "../utils/apiUtils";
import { useNavigate } from "react-router-dom";
import LoadingSpinner from "../components/LoadingSpinner";
import toast from "react-hot-toast";
import validator from "validator";

function Register() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [repeatPassword, setRepeatPassword] = useState("");

  const disableSubmit = useMemo(() => {
    return (
      username.length === 0 ||
      email.length === 0 ||
      password.length === 0 ||
      repeatPassword.length === 0 ||
      password !== repeatPassword
    );
  }, [username, email, password, repeatPassword]);

  const passwordMatch = useMemo(() => password === repeatPassword, [password, repeatPassword]);

  const validateUsername = (username: string) => {
    return username.length >= 3 && username.length <= 100;
  };

  const validateEmail = (email: string) => {
    return validator.isEmail(email);
  };

  const handleSubmit: FormEventHandler = async (e) => {
    e.preventDefault();
    if (disableSubmit) return;
    const registerData: IRegisterRequest = {
      username,
      email,
      password,
      password2: repeatPassword
    };
    setLoading(true);
    apiPost<IRegisterRequest, void>("/register", registerData)
      .then(() => {
        setLoading(false);
        toast.success("User registered successfully. Please login.");
        navigate("/login");
      })
      .catch((error) => {
        toast.error(error.message);
        setLoading(false);
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
            className="flex flex-col gap-8 lg:w-1/2 border-2 px-20 py-8 rounded bg-blue-200 bg-opacity-40"
          >
            <h1 className="text-2xl font-semibold">Register new user</h1>
            <label>
              <input
                className={`${username.length === 0 ? "" : !validateUsername(username) ? "border-red-500 outline-red-500" : "border-green-500 outline-green-500"} border-2 p-2 w-full rounded outline-2 caret-blue-500 text-xl`}
                type="text"
                placeholder="Username"
                value={username}
                min={3}
                max={100}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
              {username.length > 0 && !validateUsername(username) && (
                <span className="text-red-500 text-xs">
                  Username must be between 3 and 100 characters
                </span>
              )}
            </label>
            <label>
              <span className="text-sm">Email</span>
              <input
                className={`${email.length === 0 ? "" : !validateEmail(email) ? "border-red-500 outline-red-500" : "border-green-500 outline-green-500"} border-2 p-2 w-full rounded outline-2 caret-blue-500 text-xl`}
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
              {!validateEmail(email) && email.length > 0 && (
                <span className="text-red-500 text-xs">Invalid email address</span>
              )}
            </label>
            <label>
              <span className="text-sm">Password</span>
              <input
                className={`${password.length === 0 ? "" : password.length < 6 ? "border-red-500 outline-red-500" : "border-green-500 outline-green-500"} border-2 outline-2 p-2 w-full rounded caret-blue-500 text-xl`}
                type="password"
                placeholder="Password"
                value={password}
                minLength={6}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
              {password.length > 0 && password.length < 6 && (
                <span className="text-red-500 text-xs">Password must be at least 6 characters</span>
              )}
            </label>
            <label>
              <span className="text-sm">Repeat password</span>
              <input
                title={!passwordMatch ? "Passwords do not match" : ""}
                className={`${password.length === 0 ? "" : !passwordMatch ? "outline-red-500 border-red-500" : "outline-green-500 border-green-500"} border-2 p-2 w-full rounded outline-2 caret-blue-500 text-xl`}
                type="password"
                placeholder="Repeat password"
                value={repeatPassword}
                minLength={6}
                onChange={(e) => setRepeatPassword(e.target.value)}
                required
              />
              {!passwordMatch && (
                <span className="text-red-500 text-xs">Passwords do not match</span>
              )}
            </label>
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
              <button
                title={disableSubmit ? "Please fill out all fields" : ""}
                className={`${disableSubmit ? "grayscale-[0.5] cursor-not-allowed" : "hover:brightness-90 cursor-pointer"} border rounded bg-blue-500  text-white font-semibold p-2`}
                disabled={disableSubmit}
              >
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
