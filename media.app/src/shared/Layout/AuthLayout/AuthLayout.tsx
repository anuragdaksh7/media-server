import { Outlet } from "react-router-dom"
import styles from "./AuthLayout.module.css"
import { useState } from "react";
import { Logo } from "../../../icons";

const AuthLayout = () => {
  const [currentState, setCurrentState] = useState<"login" | "signup">("login");

  
  return (
    <div className={styles.section}>
      <div className={styles.form}>
        <div className={styles.top}>
          <div className={styles.icon}><Logo /></div>
          <h2>Media Server</h2>
          <p>ORCHESTRATOR NODE</p>
        </div>
        <div className={styles.switch_mode}>
          <div>Login</div>
        </div>
      </div>
      <Outlet />
    </div>
  )
}

export default AuthLayout