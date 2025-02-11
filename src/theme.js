import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  typography: {
    fontFamily: 'Roboto, sans-serif',
  },
  palette: {
    primary: {
      main: '#274C77', // Strong blue
    },
    secondary: {
      main: '#6096BA', // Soft blue
    },
    background: {
      default: '#E7ECEF', // Light background
    },
    text: {
      primary: '#8B8C89', // Neutral text
    },
    accent: {
      main: '#A3CEF1', // Accent blue
    },
    mode: 'light',
  },
});

export default theme;
