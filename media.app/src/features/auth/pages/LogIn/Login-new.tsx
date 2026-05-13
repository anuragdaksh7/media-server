import { useState } from "react";
import styles from "./Login-new.module.css"
import { useAppDispatch } from "../../../../app/hooks";
import { signInUser } from "../../auth.thunk";
import { useNavigate } from "react-router-dom";
import { getGoogleRedirectURL } from "../../auth.api";

const Login = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    dispatch(signInUser({ email, password }))
  }

  const handleGoogleLogin = async () => {
    const { redirect_url: redirectURL } = await getGoogleRedirectURL();
    console.log("Redirect URL:", redirectURL);
    window.location.href = redirectURL;
  }

  return (
    <div className={styles.section}>
      
    </div>
  )
}

export default Login
