import { PropsWithChildren } from "react";
import { NavLink, useNavigate } from "react-router-dom";

interface NavProps {
  loggedIn: boolean;
  onLogOut: () => void;
}

function Nav(props: Readonly<NavProps>) {
  const navigate = useNavigate();

  function logOut() {
    props.onLogOut();
    navigate("/login");
  }

  return (
    <nav className="flex w-full justify-between px-2 pt-3 bg-blue-200">
      <div className="flex">
        <CustomNavLink to="/">Search</CustomNavLink>
        <CustomNavLink to="/weather">Weather</CustomNavLink>
      </div>
      <div className="flex">
        <CustomNavLink to="/register">Register</CustomNavLink>
        <CustomNavLink
          id="login-logout-nav"
          to="/login"
          onClick={() => {
            if (props.loggedIn) logOut();
          }}
        >
          {props.loggedIn ? "Log out" : "Log in"}
        </CustomNavLink>
      </div>
    </nav>
  );
}

interface CustomNavLinkProps extends PropsWithChildren {
  to: string;
  onClick?: () => void;
  id?: string;
}

function CustomNavLink(props: Readonly<CustomNavLinkProps>) {
  const baseClassname =
    "p-2 hover:bg-white rounded-t-lg border-t border-x border-blue-200 font-semibold w-48 text-center hover:border-white transition-colors duration-500";
  return (
    <NavLink
      id={props.id}
      to={props.to}
      className={({ isActive }) => {
        return isActive ? `${baseClassname} bg-white border-white` : baseClassname;
      }}
      onClick={props.onClick}
    >
      {props.children}
    </NavLink>
  );
}

export default Nav;
