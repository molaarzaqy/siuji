import IllustrationLayout from "layouts/authentication/components/IllustrationLayout";
import SignInForm from "layouts/authentication/components/SignInForm";
import chat from "assets/images/illustrations/chat.png";

function Illustration() {
  return (
    <IllustrationLayout
      description="Masuk ke dashboard SIUJI"
      illustration={{ image: chat, title: "SIUJI", description: "Sistem Ujian Internasional" }}
    >
      <SignInForm />
    </IllustrationLayout>
  );
}

export default Illustration;
