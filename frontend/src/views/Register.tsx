import { FormEventHandler, useMemo, useState } from "react";
import type { ILoginRequest, ILoginResponse, IRegisterRequest } from "../types/auth.types";
import PageLayout from "../components/PageLayout";
import { apiPost } from "../utils/apiUtils";
import { useNavigate } from "react-router-dom";
import LoadingSpinner from "../components/LoadingSpinner";
import toast from "react-hot-toast";
import validator from "validator";

interface RegisterProps {
  logIn: (jwt_token: string) => void;
}

function Register(props: RegisterProps) {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [repeatPassword, setRepeatPassword] = useState("");

  const passwordMatch = useMemo(() => password === repeatPassword, [password, repeatPassword]);

  const disableSubmit = useMemo(() => {
    return (
      !validateUsername(username) ||
      !validateEmail(email) ||
      !validatePassword(password) ||
      !passwordMatch
    );
  }, [username, email, password, passwordMatch]);

  function validateUsername(username: string) {
    return username.length >= 3 && username.length <= 100;
  }

  function validateEmail(email: string) {
    return validator.isEmail(email);
  }

  function validatePassword(password: string) {
    return password.length >= 6;
  }

  function validateRepeatPassword(repeatPassword: string) {
    return passwordMatch && validatePassword(repeatPassword);
  }

  const handleSubmit: FormEventHandler = (e) => {
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
        toast.success("User registered successfully.");
        apiPost<ILoginRequest, ILoginResponse>("/login", { username: username, password })
          .then((data) => {
            setLoading(false);
            props.logIn(data.token);
            toast.success("Logged in successfully.");
            navigate("/");
          })
          .catch((error) => {
            setLoading(false);
            toast.error(error.message);
          });
      })
      .catch((error) => {
        toast.error(error.message);
        setLoading(false);
      });
  };

  const getInputClassName = (value: string, validator: (value: string) => boolean) => {
    const base = "border-2 p-2 w-full rounded outline-2 caret-blue-500 text-xl";
    if (value.length === 0) return base;
    return validator(value)
      ? `border-green-500 outline-green-500 ${base}`
      : `border-red-500 outline-red-500 ${base}`;
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
              <span className="text-sm">Username</span>
              <input
                className={getInputClassName(username, validateUsername)}
                type="text"
                placeholder="Username"
                value={username}
                min={3}
                max={100}
                onChange={(e) => setUsername(e.target.value)}
                required
                name="username"
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
                className={getInputClassName(email, validateEmail)}
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                name="email"
              />
              {!validateEmail(email) && email.length > 0 && (
                <span className="text-red-500 text-xs">Invalid email address</span>
              )}
            </label>
            <label>
              <span className="text-sm">Password</span>
              <input
                className={getInputClassName(password, validatePassword)}
                type="password"
                placeholder="Password"
                value={password}
                minLength={6}
                onChange={(e) => setPassword(e.target.value)}
                required
                name="password"
              />
              {password.length > 0 && password.length < 6 && (
                <span className="text-red-500 text-xs">Password must be at least 6 characters</span>
              )}
            </label>
            <label>
              <span className="text-sm">Repeat password</span>
              <input
                title={!passwordMatch ? "Passwords do not match" : ""}
                className={getInputClassName(repeatPassword, validateRepeatPassword)}
                type="password"
                placeholder="Repeat password"
                value={repeatPassword}
                minLength={6}
                onChange={(e) => setRepeatPassword(e.target.value)}
                required
                name="repeat-password"
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
                id="register-button"
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
