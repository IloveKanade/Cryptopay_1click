import { createTheme } from '@mui/material/styles';

export const geekTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#00ff88', // 霓虹绿
      light: '#66ffaa',
      dark: '#00cc6a',
      contrastText: '#000',
    },
    secondary: {
      main: '#ff0080', // 霓虹粉
      light: '#ff66b3',
      dark: '#cc0066',
      contrastText: '#fff',
    },
    background: {
      default: '#0a0a0a', // 深黑色背景
      paper: '#1a1a1a', // 稍亮的卡片背景
    },
    text: {
      primary: '#ffffff',
      secondary: '#b0b0b0',
    },
    error: {
      main: '#ff4444', // 霓虹红
    },
    warning: {
      main: '#ffaa00', // 霓虹橙
    },
    info: {
      main: '#0088ff', // 霓虹蓝
    },
    success: {
      main: '#00ff88', // 霓虹绿
    },
  },
  typography: {
    fontFamily: '"Orbitron", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontWeight: 700,
      textTransform: 'uppercase',
      letterSpacing: '0.1em',
    },
    h2: {
      fontWeight: 600,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    h3: {
      fontWeight: 600,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    h4: {
      fontWeight: 500,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    h5: {
      fontWeight: 500,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    h6: {
      fontWeight: 500,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    button: {
      textTransform: 'uppercase',
      fontWeight: 600,
      letterSpacing: '0.1em',
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 25,
          padding: '12px 24px',
          boxShadow: '0 4px 15px rgba(0, 255, 136, 0.3)',
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'translateY(-2px)',
            boxShadow: '0 6px 20px rgba(0, 255, 136, 0.5)',
          },
        },
        contained: {
          background: 'linear-gradient(45deg, #00ff88 30%, #00cc6a 90%)',
          '&:hover': {
            background: 'linear-gradient(45deg, #00cc6a 30%, #00ff88 90%)',
          },
        },
        outlined: {
          borderColor: '#00ff88',
          color: '#00ff88',
          '&:hover': {
            borderColor: '#00cc6a',
            backgroundColor: 'rgba(0, 255, 136, 0.1)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          background: 'linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%)',
          border: '1px solid rgba(0, 255, 136, 0.2)',
          boxShadow: '0 8px 32px rgba(0, 255, 136, 0.1)',
          backdropFilter: 'blur(10px)',
          '&:hover': {
            borderColor: 'rgba(0, 255, 136, 0.4)',
            boxShadow: '0 12px 40px rgba(0, 255, 136, 0.2)',
            transform: 'translateY(-4px)',
            transition: 'all 0.3s ease',
          },
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          background: 'linear-gradient(90deg, #0a0a0a 0%, #1a1a1a 100%)',
          borderBottom: '1px solid rgba(0, 255, 136, 0.3)',
          backdropFilter: 'blur(10px)',
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            '& fieldset': {
              borderColor: 'rgba(0, 255, 136, 0.3)',
            },
            '&:hover fieldset': {
              borderColor: 'rgba(0, 255, 136, 0.5)',
            },
            '&.Mui-focused fieldset': {
              borderColor: '#00ff88',
            },
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          background: 'linear-gradient(45deg, #00ff88 30%, #00cc6a 90%)',
          color: '#000',
          fontWeight: 600,
          '&.MuiChip-colorDefault': {
            background: 'linear-gradient(45deg, #666 30%, #888 90%)',
            color: '#fff',
          },
        },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: {
          background: 'linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%)',
          border: '1px solid rgba(0, 255, 136, 0.3)',
          boxShadow: '0 20px 60px rgba(0, 255, 136, 0.2)',
        },
      },
    },
  },
});

export const lightGeekTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#00ff88', // 霓虹绿
      light: '#66ffaa',
      dark: '#00cc6a',
      contrastText: '#000',
    },
    secondary: {
      main: '#ff0080', // 霓虹粉
      light: '#ff66b3',
      dark: '#cc0066',
      contrastText: '#fff',
    },
    background: {
      default: '#f8f9fa',
      paper: '#ffffff',
    },
    text: {
      primary: '#1a1a1a',
      secondary: '#666666',
    },
  },
  typography: {
    fontFamily: '"Orbitron", "Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontWeight: 700,
      textTransform: 'uppercase',
      letterSpacing: '0.1em',
    },
    h2: {
      fontWeight: 600,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    h3: {
      fontWeight: 600,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    h4: {
      fontWeight: 500,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    h5: {
      fontWeight: 500,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    h6: {
      fontWeight: 500,
      textTransform: 'uppercase',
      letterSpacing: '0.05em',
    },
    button: {
      textTransform: 'uppercase',
      fontWeight: 600,
      letterSpacing: '0.1em',
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 25,
          padding: '12px 24px',
          boxShadow: '0 4px 15px rgba(0, 255, 136, 0.3)',
          transition: 'all 0.3s ease',
          '&:hover': {
            transform: 'translateY(-2px)',
            boxShadow: '0 6px 20px rgba(0, 255, 136, 0.5)',
          },
        },
        contained: {
          background: 'linear-gradient(45deg, #00ff88 30%, #00cc6a 90%)',
          '&:hover': {
            background: 'linear-gradient(45deg, #00cc6a 30%, #00ff88 90%)',
          },
        },
        outlined: {
          borderColor: '#00ff88',
          color: '#00ff88',
          '&:hover': {
            borderColor: '#00cc6a',
            backgroundColor: 'rgba(0, 255, 136, 0.1)',
          },
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          border: '1px solid rgba(0, 255, 136, 0.2)',
          boxShadow: '0 8px 32px rgba(0, 255, 136, 0.1)',
          '&:hover': {
            borderColor: 'rgba(0, 255, 136, 0.4)',
            boxShadow: '0 12px 40px rgba(0, 255, 136, 0.2)',
            transform: 'translateY(-4px)',
            transition: 'all 0.3s ease',
          },
        },
      },
    },
  },
});
