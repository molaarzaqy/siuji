import React, { forwardRef } from "react";
import PropTypes from "prop-types";
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";

const CertificateTemplate = forwardRef(({ participantName, periodTitle, score, dateStr }, ref) => {
  return (
    <div style={{ overflow: "hidden", position: "absolute", width: 0, height: 0, opacity: 0, zIndex: -9999 }}>
      {/* 
        Fixed A4 Size landscape: 1123px x 794px 
        Using standard CSS for html2canvas compatibility
      */}
      <div 
        ref={ref} 
        style={{
          width: "1123px",
          height: "794px",
          background: "linear-gradient(135deg, #ffffff, #f0f4f8)",
          position: "relative",
          padding: "40px",
          boxSizing: "border-box",
          fontFamily: "Arial, sans-serif"
        }}
      >
        <div style={{
          width: "100%",
          height: "100%",
          border: "8px solid #344767", // dark theme color
          borderRadius: "16px",
          position: "relative",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center"
        }}>
          {/* Inner border */}
          <div style={{
            position: "absolute",
            top: "10px", bottom: "10px", left: "10px", right: "10px",
            border: "2px solid #344767",
            borderRadius: "8px"
          }}></div>

          <SoftBox textAlign="center" zIndex={2}>
            <SoftTypography variant="h1" fontWeight="bold" sx={{ color: "#344767", fontSize: "48px", textTransform: "uppercase", letterSpacing: "4px" }}>
              Sertifikat Kompetensi
            </SoftTypography>
            
            <SoftTypography variant="h5" sx={{ color: "#8392ab", marginTop: "16px", fontStyle: "italic" }}>
              Diberikan kepada:
            </SoftTypography>

            <SoftTypography variant="h2" fontWeight="bold" sx={{ color: "#17c1e8", marginTop: "24px", fontSize: "42px", textTransform: "capitalize" }}>
              {participantName}
            </SoftTypography>

            <SoftTypography variant="body1" sx={{ color: "#67748e", marginTop: "24px", maxWidth: "600px", marginX: "auto", fontSize: "20px" }}>
              Telah berhasil menyelesaikan ujian <strong>{periodTitle}</strong> yang diselenggarakan pada platform SIUJI dengan hasil yang memuaskan.
            </SoftTypography>

            <SoftBox mt={6} display="inline-flex" alignItems="center" justifyContent="center" 
              sx={{ 
                background: "linear-gradient(135deg, #43a047, #66bb6a)", 
                borderRadius: "50%", 
                width: "140px", height: "140px",
                border: "4px solid #fff",
                boxShadow: "0 8px 32px rgba(67, 160, 71, 0.3)"
              }}
            >
              <SoftBox>
                <SoftTypography variant="h3" fontWeight="bold" color="white">{score}</SoftTypography>
                <SoftTypography variant="caption" color="white" display="block">TOEFL Score</SoftTypography>
              </SoftBox>
            </SoftBox>

            <SoftBox mt={6} display="flex" justifyContent="space-around" width="100%">
              <SoftBox>
                <SoftTypography variant="body2" color="text">Diterbitkan pada:</SoftTypography>
                <SoftTypography variant="h6" fontWeight="bold">{dateStr}</SoftTypography>
              </SoftBox>
              <SoftBox>
                <SoftTypography variant="body2" color="text">Penyelenggara</SoftTypography>
                <SoftTypography variant="h6" fontWeight="bold">SIUJI Platform</SoftTypography>
              </SoftBox>
            </SoftBox>

          </SoftBox>
        </div>
      </div>
    </div>
  );
});

CertificateTemplate.propTypes = {
  participantName: PropTypes.string.isRequired,
  periodTitle: PropTypes.string.isRequired,
  score: PropTypes.number.isRequired,
  dateStr: PropTypes.string.isRequired,
};

export default CertificateTemplate;
