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

import { forwardRef } from "react";

// prop-types is a library for typechecking of props
import PropTypes from "prop-types";

// react-select components
import Select from "react-select";

// Soft UI Dashboard PRO React base styles
import colors from "assets/theme/base/colors";

// Custom styles for SoftSelect
import styles from "components/SoftSelect/styles";

const SoftSelect = forwardRef(({
  size = "medium",
  error = false,
  success = false,
  isClearable = true,
  noOptionsMessage = () => "Tidak ada pilihan",
  loadingMessage = () => "Memuat pilihan...",
  ...rest
}, ref) => {
  const { light } = colors;

  return (
    <Select
      {...rest}
      ref={ref}
      isClearable={isClearable}
      noOptionsMessage={noOptionsMessage}
      loadingMessage={loadingMessage}
      styles={styles(size, error, success)}
      menuPortalTarget={typeof document !== "undefined" ? document.body : null}
      menuPosition="fixed"
      theme={(theme) => ({
        ...theme,
        colors: {
          ...theme.colors,
          primary25: light.main,
          primary: light.main,
        },
      })}
    />
  );
});

// Typechecking props for the SoftSelect
SoftSelect.propTypes = {
  size: PropTypes.oneOf(["small", "medium", "large"]),
  error: PropTypes.bool,
  success: PropTypes.bool,
  isClearable: PropTypes.bool,
  noOptionsMessage: PropTypes.func,
  loadingMessage: PropTypes.func,
};

export default SoftSelect;
