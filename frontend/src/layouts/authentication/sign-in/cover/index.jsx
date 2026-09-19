import CoverLayout from "layouts/authentication/components/CoverLayout";
import SignInForm from "layouts/authentication/components/SignInForm";
import curved9 from "assets/images/curved-images/curved9.jpg";

function Cover() {
  return <CoverLayout description="Masuk ke dashboard SIUJI" image={curved9}><SignInForm /></CoverLayout>;
}

export default Cover;
