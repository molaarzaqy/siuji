/**
=========================================================
* Soft UI Dashboard PRO React - v4.0.3
=========================================================

* Product Page: https://www.creative-tim.com/product/soft-ui-dashboard-pro-react
* Copyright 2024 Creative Tim (https://www.creative-tim.com)

Coded by www.creative-tim.com

 =========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
*/

import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "context/auth/AuthContext";

// prop-types is a library for typechecking of props
import PropTypes from "prop-types";

// @mui material components
import Grid from "@mui/material/Grid";
import AppBar from "@mui/material/AppBar";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Icon from "@mui/material/Icon";

// Soft UI Dashboard PRO React components
import SoftBox from "components/SoftBox";

// Soft UI Dashboard PRO React base styles
import breakpoints from "assets/theme/base/breakpoints";

// Soft UI Dashboard PRO React example components
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";

function BaseLayout({ stickyNavbar = false, children }) {
  const [tabsOrientation, setTabsOrientation] = useState("horizontal");
  const [tabValue, setTabValue] = useState(0);
  const navigate = useNavigate();
  const { user } = useAuth();
  const role = user?.role || "participant";

  useEffect(() => {
    // A function that sets the orientation state of the tabs.
    function handleTabsOrientation() {
      return window.innerWidth < breakpoints.values.sm
        ? setTabsOrientation("vertical")
        : setTabsOrientation("horizontal");
    }

    /** 
     The event listener that's calling the handleTabsOrientation function when resizing the window.
    */
    window.addEventListener("resize", handleTabsOrientation);

    // Call the handleTabsOrientation function to set the state with the initial value.
    handleTabsOrientation();

    // Remove event listener on cleanup
    return () => window.removeEventListener("resize", handleTabsOrientation);
  }, [tabsOrientation]);

  const handleSetTabValue = (event, newValue) => {
    setTabValue(newValue);
    
    // Quick Navigation Logic based on Tab Index
    if (role === "admin") {
      if (newValue === 1) navigate("/dashboards/default");
      if (newValue === 2) navigate("/admin/periods");
      if (newValue === 3) navigate("/admin/sections");
    } else {
      if (newValue === 1) navigate("/participant/periods");
      // For certificates & help, we can keep them mock for now
    }
  };

  return (
    <DashboardLayout>
      <DashboardNavbar absolute={!stickyNavbar} isMini />
      <SoftBox mt={stickyNavbar ? 3 : 10}>
        {children}
      </SoftBox>
      <Footer />
    </DashboardLayout>
  );
}


// Typechecking props for BaseLayout
BaseLayout.propTypes = {
  stickyNavbar: PropTypes.bool,
  children: PropTypes.node.isRequired,
};

export default BaseLayout;
