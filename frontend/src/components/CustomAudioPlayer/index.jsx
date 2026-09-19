import { useState, useRef, useEffect } from "react";
import PropTypes from "prop-types";
import SoftBox from "components/SoftBox";
import SoftButton from "components/SoftButton";
import SoftTypography from "components/SoftTypography";
import Icon from "@mui/material/Icon";
import LinearProgress from "@mui/material/LinearProgress";

function CustomAudioPlayer({ src }) {
  const audioRef = useRef(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const [duration, setDuration] = useState(0);
  const [currentTime, setCurrentTime] = useState(0);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const updateProgress = () => {
      setCurrentTime(audio.currentTime);
      setProgress((audio.currentTime / audio.duration) * 100 || 0);
    };

    const handleLoadedMetadata = () => {
      setDuration(audio.duration);
    };

    const handleEnded = () => {
      setIsPlaying(false);
      setProgress(100);
    };

    audio.addEventListener("timeupdate", updateProgress);
    audio.addEventListener("loadedmetadata", handleLoadedMetadata);
    audio.addEventListener("ended", handleEnded);

    return () => {
      audio.removeEventListener("timeupdate", updateProgress);
      audio.removeEventListener("loadedmetadata", handleLoadedMetadata);
      audio.removeEventListener("ended", handleEnded);
    };
  }, []);

  const togglePlay = () => {
    const audio = audioRef.current;
    if (!audio) return;

    if (isPlaying) {
      audio.pause();
    } else {
      audio.play().catch(console.error);
    }
    setIsPlaying(!isPlaying);
  };

  const formatTime = (time) => {
    if (isNaN(time)) return "0:00";
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60);
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  };

  const handleSeek = (event) => {
    const audio = audioRef.current;
    if (!audio) return;
    
    // Calculate new time based on click position
    const bounds = event.currentTarget.getBoundingClientRect();
    const percent = (event.clientX - bounds.left) / bounds.width;
    const newTime = percent * audio.duration;
    
    audio.currentTime = newTime;
    setCurrentTime(newTime);
    setProgress(percent * 100);
  };

  return (
    <SoftBox
      display="flex"
      alignItems="center"
      p={2}
      bgColor="grey.100"
      borderRadius="lg"
      gap={2}
      sx={{ boxShadow: "inset 0 0 0 1px rgba(0,0,0,0.05)" }}
    >
      <audio ref={audioRef} src={src} preload="metadata" />
      
      <SoftButton
        variant="gradient"
        color="info"
        iconOnly
        circular
        onClick={togglePlay}
        sx={{ minWidth: 48, width: 48, height: 48 }}
      >
        <Icon fontSize="medium">{isPlaying ? "pause" : "play_arrow"}</Icon>
      </SoftButton>

      <SoftBox flex={1} display="flex" flexDirection="column" gap={0.5}>
        <SoftBox display="flex" justifyContent="space-between">
          <SoftTypography variant="caption" fontWeight="bold" color="text">
            {formatTime(currentTime)}
          </SoftTypography>
          <SoftTypography variant="caption" fontWeight="bold" color="text">
            {formatTime(duration)}
          </SoftTypography>
        </SoftBox>
        
        <SoftBox 
          onClick={handleSeek}
          sx={{ 
            cursor: "pointer", 
            py: 1,
            "&:hover .MuiLinearProgress-root": { height: 6 }
          }}
        >
          <LinearProgress
            variant="determinate"
            value={progress}
            color="info"
            sx={{
              height: 4,
              borderRadius: 2,
              transition: "height 0.2s",
              backgroundColor: "grey.300"
            }}
          />
        </SoftBox>
      </SoftBox>
    </SoftBox>
  );
}

CustomAudioPlayer.propTypes = {
  src: PropTypes.string.isRequired,
};

export default CustomAudioPlayer;
