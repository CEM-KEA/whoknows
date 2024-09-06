import { PropsWithChildren } from "react";
import { NavLink } from "react-router-dom";

interface NavProps {
    loggedIn: boolean;
}

function Nav(props: NavProps) {
    return (
        <nav className="flex w-full justify-between px-2 pt-3 bg-gray-200">
            <div className="flex">
                <CustomNavLink to="/">Search</CustomNavLink>
                <CustomNavLink to="/weather">Weather</CustomNavLink>
            </div>
            <div className="flex">
                <CustomNavLink to="/register">Register</CustomNavLink>
                <CustomNavLink to="/login">{props.loggedIn ? "Log out" : "Login"}</CustomNavLink>
            </div>
        </nav>
    )
}

interface CustomNavLinkProps extends PropsWithChildren {
    to: string;
}


function CustomNavLink(props: CustomNavLinkProps) {
    const baseClassname = "p-2 hover:bg-white rounded-t-lg border font-semibold w-48 text-center hover:border-white transition-colors duration-300";
    return (
        <NavLink 
            to={props.to} 
            className={
                ({ isActive }) => {
                    return isActive ? `${baseClassname} bg-white border-white` : baseClassname
            }}
        >
            {props.children}
        </NavLink>
    )
}

export default Nav;