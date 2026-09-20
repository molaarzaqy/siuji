// SIUJI layouts
import Default from "layouts/dashboards/default";
import Periods from "layouts/admin/periods";
import PeriodDetail from "layouts/admin/period-detail";
import Sections from "layouts/admin/sections";
import Participants from "layouts/admin/participants";
import Users from "layouts/admin/users";
import ScoreConversions from "layouts/admin/score-conversions";
import ParticipantDashboard from "layouts/participant/dashboard";
import ParticipantPeriods from "layouts/participant/periods";
import ParticipantExam from "layouts/participant/exam";
import ParticipantResult from "layouts/participant/results";

// Account / Profile Settings
import Settings from "layouts/pages/account/settings";

// Authentication layouts
import SignInBasic from "layouts/authentication/sign-in/basic";
import ResetBasic from "layouts/authentication/reset-password/basic";
import UpdatePassword from "layouts/authentication/reset-password/update";
import VerificationBasic from "layouts/authentication/2-step-verification/basic";
import Error404 from "layouts/authentication/error/404";
import Error500 from "layouts/authentication/error/500";

// React icons
import Shop from "examples/Icons/Shop";
import Office from "examples/Icons/Office";
import Document from "examples/Icons/Document";
import CreditCard from "examples/Icons/CreditCard";
import CustomerSupport from "examples/Icons/CustomerSupport";
import SettingsIcon from "examples/Icons/Settings";

const routes = [
  // ── Admin: Dashboard ───────────────────────────────────
  {
    type: "collapse",
    name: "Dashboard",
    key: "dashboards",
    roles: ["admin"],
    icon: <Shop size="12px" />,
    noCollapse: true,
    route: "/dashboards/default",
    component: <Default />,
  },

  // ── Admin: Administrasi Ujian ──────────────────────────
  {
    type: "collapse",
    name: "Administrasi Ujian",
    key: "administration",
    roles: ["admin"],
    icon: <Office size="12px" />,
    collapse: [
      {
        name: "Periode Ujian",
        key: "periods",
        route: "/admin/periods",
        component: <Periods />,
        roles: ["admin"],
      },
      {
        name: "Peserta Ujian",
        key: "participants",
        route: "/admin/participants",
        component: <Participants />,
        roles: ["admin"],
      },
      {
        name: "User",
        key: "admin-users",
        route: "/admin/users",
        component: <Users />,
        roles: ["admin"],
      },
    ],
  },

  // ── Admin: Bank Soal ───────────────────────────────────
  {
    type: "collapse",
    name: "Bank Soal",
    key: "question-bank",
    roles: ["admin"],
    icon: <Document size="12px" />,
    collapse: [
      {
        name: "Section",
        key: "sections",
        route: "/admin/sections",
        component: <Sections />,
        roles: ["admin"],
      },
    ],
  },

  // ── Admin: Scoring ─────────────────────────────────────
  {
    type: "collapse",
    name: "Scoring",
    key: "scoring",
    roles: ["admin"],
    icon: <CreditCard size="12px" />,
    collapse: [
      {
        name: "Konversi Nilai",
        key: "score-conversions",
        route: "/admin/score-conversions",
        component: <ScoreConversions />,
        roles: ["admin"],
      },
    ],
  },

  // ── Participant: Dashboard ──────────────────────────────
  {
    type: "collapse",
    name: "Dashboard",
    key: "participant-dashboard",
    roles: ["participant"],
    icon: <Shop size="12px" />,
    noCollapse: true,
    route: "/participant/dashboard",
    component: <ParticipantDashboard />,
  },

  // ── Participant: Manajemen Ujian ─────────────────────────
  {
    type: "collapse",
    name: "Ujian Saya",
    key: "participant-periods",
    roles: ["participant"],
    icon: <Document size="12px" />,
    noCollapse: true,
    route: "/participant/periods",
    component: <ParticipantPeriods />,
  },

  // ── Account / Profile ──────────────────────────────────
  {
    type: "collapse",
    name: "Pengaturan Akun",
    key: "profile",
    roles: ["admin", "participant"],
    icon: <SettingsIcon size="12px" />,
    noCollapse: true,
    route: "/profile/settings",
    component: <Settings />,
  },

  // ── Hidden routes (tidak tampil di sidebar) ────────────
  {
    type: "hidden",
    key: "period-detail",
    route: "/admin/periods/:id",
    component: <PeriodDetail />,
    roles: ["admin"],
  },
  {
    type: "hidden",
    key: "participant-exam",
    route: "/participant/exam/:id",
    component: <ParticipantExam />,
    roles: ["participant"],
  },
  {
    type: "hidden",
    key: "participant-results",
    route: "/participant/results/:id",
    component: <ParticipantResult />,
    roles: ["participant"],
  },

  // ── Authentication (hidden from sidebar) ───────────────
  {
    type: "hidden",
    key: "sign-in",
    route: "/authentication/sign-in/basic",
    component: <SignInBasic />,
  },
  {
    type: "hidden",
    key: "reset-password",
    route: "/authentication/reset-password/basic",
    component: <ResetBasic />,
  },
  {
    type: "hidden",
    key: "update-password",
    route: "/authentication/reset-password/update",
    component: <UpdatePassword />,
  },
  {
    type: "hidden",
    key: "verification",
    route: "/authentication/verification/basic",
    component: <VerificationBasic />,
  },
  {
    type: "hidden",
    key: "error-404",
    route: "/authentication/error/404",
    component: <Error404 />,
  },
  {
    type: "hidden",
    key: "error-500",
    route: "/authentication/error/500",
    component: <Error500 />,
  },
];

export default routes;
