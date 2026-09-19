import { useState, useEffect, useMemo } from "react";

// react-router components
import { Routes, Route, Navigate, useLocation } from "react-router-dom";

// @mui material components
import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import Icon from "@mui/material/Icon";

//  PRO React components
import SoftBox from "components/SoftBox";

//  PRO React example components
import Sidenav from "examples/Sidenav";
import Configurator from "examples/Configurator";
import PrivateRoute, { PublicRoute, RoleRoute } from "examples/PrivateRoute";
import LandingPage from "layouts/landing";
import { useAuth } from "context/auth/AuthContext";

//  PRO React themes
import theme from "assets/theme";
import themeRTL from "assets/theme/theme-rtl";

// RTL plugins
import rtlPlugin from "stylis-plugin-rtl";
import { CacheProvider } from "@emotion/react";
import createCache from "@emotion/cache";

//  PRO React routes
import routes from "routes";

//  PRO React contexts
import { useSoftUIController, setMiniSidenav, setOpenConfigurator } from "context";

// Images
import bizbeamLogo from "assets/images/bizbeam-logo (1).png";

export default function App() {
  const { user } = useAuth();
  const [controller, dispatch] = useSoftUIController();
  const { miniSidenav, direction, layout, openConfigurator, sidenavColor } = controller;
  const [onMouseEnter, setOnMouseEnter] = useState(false);
  const [rtlCache, setRtlCache] = useState(null);
  const { pathname } = useLocation();
  const isAuthenticationPage = pathname.startsWith("/authentication/");
  const filterRoutesByRole = (items) => items
    .filter((route) => !route.roles || !user || route.roles.includes(user.role))
    .map((route) => {
      if (!route.collapse) return route;
      const collapse = filterRoutesByRole(route.collapse);
      return collapse.length || route.route ? { ...route, collapse } : null;
    })
    .filter(Boolean);
  const visibleRoutes = filterRoutesByRole(routes);

  // Cache for the rtl
  useMemo(() => {
    const cacheRtl = createCache({
      key: "rtl",
      stylisPlugins: [rtlPlugin],
    });

    setRtlCache(cacheRtl);
  }, []);

  // Open sidenav when mouse enter on mini sidenav
  const handleOnMouseEnter = () => {
    if (miniSidenav && !onMouseEnter) {
      setMiniSidenav(dispatch, false);
      setOnMouseEnter(true);
    }
  };

  // Close sidenav when mouse leave mini sidenav
  const handleOnMouseLeave = () => {
    if (onMouseEnter) {
      setMiniSidenav(dispatch, true);
      setOnMouseEnter(false);
    }
  };

  // Change the openConfigurator state
  const handleConfiguratorOpen = () => setOpenConfigurator(dispatch, !openConfigurator);

  // Setting the dir attribute for the body element
  useEffect(() => {
    document.body.setAttribute("dir", direction);
  }, [direction]);

  // Setting page scroll to 0 when changing the route
  useEffect(() => {
    document.documentElement.scrollTop = 0;
    document.scrollingElement.scrollTop = 0;
  }, [pathname]);


  const getRoutes = (allRoutes) =>
    allRoutes.map((route) => {
      if (route.collapse) {
        return getRoutes(route.collapse);
      }

      if (route.route) {

        // Authentication routes
        if (route.route.includes("/authentication/") && route.route !== "/authentication/reset-password/update") {
          const isSignInRoute = route.route.startsWith("/authentication/sign-in/");
          return (
            <Route
              exact
              path={route.route}
              element={isSignInRoute ? <PublicRoute>{route.component}</PublicRoute> : route.component}
              key={route.key}
            />
          );
        }

        // Protected routes (including hidden type routes)
        const protectedComponent = route.roles
          ? <RoleRoute roles={route.roles}>{route.component}</RoleRoute>
          : route.component;
        return (
          <Route
            exact
            path={route.route}
            element={<PrivateRoute>{protectedComponent}</PrivateRoute>}
            key={route.key}
          />
        );
      }

      return null;
    });

  const configsButton = (
    <SoftBox
      display="flex"
      justifyContent="center"
      alignItems="center"
      width="3.5rem"
      height="3.5rem"
      bgColor="white"
      shadow="sm"
      borderRadius="50%"
      position="fixed"
      right="2rem"
      bottom="2rem"
      zIndex={99}
      color="dark"
      sx={{ cursor: "pointer" }}
      onClick={handleConfiguratorOpen}
    >
      <Icon fontSize="default" color="inherit">
        settings
      </Icon>
    </SoftBox>
  );

  const isLandingPage = pathname === "/";

  return direction === "rtl" ? (
    <CacheProvider value={rtlCache}>
      <ThemeProvider theme={themeRTL}>
        <CssBaseline />
        {layout === "dashboard" && !isAuthenticationPage && !isLandingPage && (
          <>
            <Sidenav
              color={sidenavColor}
              brand={bizbeamLogo}
              routes={visibleRoutes}
              onMouseEnter={handleOnMouseEnter}
              onMouseLeave={handleOnMouseLeave}
            />
            <Configurator />
            {configsButton}
          </>
        )}
        <Routes>
          <Route path="/" element={<PublicRoute><LandingPage /></PublicRoute>} />
          {getRoutes(routes)}
          <Route path="*" element={<Navigate to={user?.role === "participant" ? "/participant/periods" : "/dashboards/default"} />} />
        </Routes>
      </ThemeProvider>
    </CacheProvider>
  ) : (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {layout === "dashboard" && !isAuthenticationPage && !isLandingPage && (
        <>
          <Sidenav
            color={sidenavColor}
            brand={bizbeamLogo}
          routes={visibleRoutes}
            onMouseEnter={handleOnMouseEnter}
            onMouseLeave={handleOnMouseLeave}
          />
          <Configurator />
          {configsButton}
        </>
      )}
      <Routes>
          <Route path="/" element={<PublicRoute><LandingPage /></PublicRoute>} />
          {getRoutes(routes)}
        <Route path="*" element={<Navigate to={user?.role === "participant" ? "/participant/periods" : "/dashboards/default"} />} />
      </Routes>
    </ThemeProvider>
  );
}

