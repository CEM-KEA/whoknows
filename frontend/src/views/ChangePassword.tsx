import { FormEventHandler, useMemo, useState } from "react";
import PageLayout from "../components/PageLayout";
import { IChangePasswordRequest } from "../types/auth.types";
import LoadingSpinner from "../components/LoadingSpinner";
import { apiPost } from "../utils/apiUtils";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";

interface ChangePasswordProps {
  loggedIn: boolean;
  logOut: () => void;
}

function ChangePassword(props: ChangePasswordProps) {
  const navigate = useNavigate();
  const [username, setUsername] = useState<string>("");
  const [oldPassword, setOldPassword] = useState<string>("");
  const [newPassword, setNewPassword] = useState<string>("");
  const [repeatNewPassword, setRepeatNewPassword] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const passwordMatch = useMemo(
    () => newPassword === repeatNewPassword,
    [newPassword, repeatNewPassword]
  );

  const disableSubmit = useMemo(
    () => username === "" || oldPassword === "" || !validatePassword(newPassword) || !passwordMatch,
    [username, oldPassword, newPassword, passwordMatch]
  );
  function validatePassword(password: string) {
    return password.length >= 6;
  }

  function validateRepeatPassword(repeatPassword: string) {
    return passwordMatch && validatePassword(repeatPassword);
  }

  const handleSubmit: FormEventHandler = (e) => {
    e.preventDefault();
    if (disableSubmit) return;
    const changePasswordData: IChangePasswordRequest = {
      username: username,
      old_password: oldPassword,
      new_password: newPassword,
      repeat_new_password: repeatNewPassword
    };
    setLoading(true);
    apiPost<IChangePasswordRequest, void>("/change-password", changePasswordData)
      .then(() => {
        toast.success("Password changed successfully");
        if (props.loggedIn) props.logOut();
        navigate("/login");
      })
      .catch((e) => {
        setLoading(false);
        toast.error(e.message);
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
            <h1 className="text-2xl font-semibold">Change password</h1>
            <label>
              <span className="text-sm">Username</span>
              <input
                className={getInputClassName(username, (name) => name !== "")}
                type="text"
                placeholder="Username"
                value={username}
                min={3}
                max={100}
                onChange={(e) => setUsername(e.target.value)}
                required
                name="username"
              />
            </label>
            <label>
              <span className="text-sm">Old Password</span>
              <input
                className={getInputClassName(oldPassword, (pass) => pass !== "")}
                type="password"
                placeholder="Old Password"
                value={oldPassword}
                minLength={6}
                onChange={(e) => setOldPassword(e.target.value)}
                required
                name="oldpassword"
              />
            </label>
            <label>
              <span className="text-sm">New Password</span>
              <input
                className={getInputClassName(newPassword, validatePassword)}
                type="password"
                placeholder="New Password"
                value={newPassword}
                minLength={6}
                onChange={(e) => setNewPassword(e.target.value)}
                required
                name="newpassword"
              />
              {newPassword.length > 0 && newPassword.length < 6 && (
                <span className="text-red-500 text-xs">Password must be at least 6 characters</span>
              )}
            </label>
            <label>
              <span className="text-sm">Repeat new password</span>
              <input
                title={!passwordMatch ? "Passwords do not match" : ""}
                className={getInputClassName(repeatNewPassword, validateRepeatPassword)}
                type="password"
                placeholder="Repeat new password"
                value={repeatNewPassword}
                minLength={6}
                onChange={(e) => setRepeatNewPassword(e.target.value)}
                required
                name="repeat-newpassword"
              />
              {!passwordMatch && (
                <span className="text-red-500 text-xs">Passwords do not match</span>
              )}
            </label>
            <div className="flex gap-1 w-full justify-between mt-8">
              <button
                className="border rounded bg-blue-50 text-blue-500 hover:bg-blue-100 font-semibold p-2"
                onClick={() => {}}
                type="button"
              >
                Clear
              </button>
              <button
                title={disableSubmit ? "Please fill out all fields" : ""}
                className={`${disableSubmit ? "grayscale-[0.5] cursor-not-allowed" : "hover:brightness-90 cursor-pointer"} border rounded bg-blue-500  text-white font-semibold p-2`}
                disabled={disableSubmit}
                id="change-password-button"
              >
                Submit
              </button>
            </div>
          </form>
        )}
      </div>
    </PageLayout>
  );
}

export default ChangePassword;
