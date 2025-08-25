import React, { useState } from 'react';
import {
  Box,
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Link,
  Alert,
} from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const RegisterPage: React.FC = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError('两次输入的密码不一致');
      return;
    }

    if (password.length < 6) {
      setError('密码长度至少6位');
      return;
    }

    setLoading(true);

    try {
      await register(name, email, password);
      navigate('/login', { state: { message: '注册成功，请登录' } });
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ 
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #16213e 100%)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      py: 4
    }}>
      <Container component="main" maxWidth="sm">
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
          }}
        >
          <Paper
            elevation={0}
            sx={{
              padding: 4,
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              width: '100%',
              background: 'rgba(26, 26, 46, 0.9)',
              backdropFilter: 'blur(20px)',
              border: '1px solid rgba(0, 255, 136, 0.3)',
              borderRadius: 3,
              boxShadow: '0 20px 40px rgba(0, 0, 0, 0.5)'
            }}
          >
            <Typography 
              component="h1" 
              variant="h5" 
              gutterBottom
              sx={{ 
                color: '#ffffff',
                fontWeight: 'bold',
                mb: 3
              }}
            >
              注册
            </Typography>

            {error && (
              <Alert 
                severity="error" 
                sx={{ 
                  width: '100%', 
                  mb: 3,
                  background: 'rgba(255, 68, 68, 0.1)',
                  border: '1px solid rgba(255, 68, 68, 0.3)',
                  color: '#ff4444'
                }}
              >
                {error}
              </Alert>
            )}

            <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1, width: '100%' }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="name"
                label="姓名"
                name="name"
                autoComplete="name"
                autoFocus
                value={name}
                onChange={(e) => setName(e.target.value)}
                sx={{
                  mb: 2,
                  '& .MuiOutlinedInput-root': {
                    color: '#ffffff',
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
                  '& .MuiInputLabel-root': {
                    color: '#b0b0b0',
                    '&.Mui-focused': {
                      color: '#00ff88',
                    },
                  },
                }}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="邮箱地址"
                name="email"
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                sx={{
                  mb: 2,
                  '& .MuiOutlinedInput-root': {
                    color: '#ffffff',
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
                  '& .MuiInputLabel-root': {
                    color: '#b0b0b0',
                    '&.Mui-focused': {
                      color: '#00ff88',
                    },
                  },
                }}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="密码"
                type="password"
                id="password"
                autoComplete="new-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                sx={{
                  mb: 2,
                  '& .MuiOutlinedInput-root': {
                    color: '#ffffff',
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
                  '& .MuiInputLabel-root': {
                    color: '#b0b0b0',
                    '&.Mui-focused': {
                      color: '#00ff88',
                    },
                  },
                }}
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="confirmPassword"
                label="确认密码"
                type="password"
                id="confirmPassword"
                autoComplete="new-password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                sx={{
                  mb: 3,
                  '& .MuiOutlinedInput-root': {
                    color: '#ffffff',
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
                  '& .MuiInputLabel-root': {
                    color: '#b0b0b0',
                    '&.Mui-focused': {
                      color: '#00ff88',
                    },
                  },
                }}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ 
                  mt: 3, 
                  mb: 3,
                  background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                  color: '#000',
                  fontWeight: 'bold',
                  py: 1.5,
                  borderRadius: 2,
                  textTransform: 'none',
                  fontSize: '1.1rem',
                  boxShadow: '0 8px 32px rgba(0, 255, 136, 0.3)',
                  '&:hover': {
                    background: 'linear-gradient(45deg, #00cc6a, #0099cc)',
                    boxShadow: '0 12px 40px rgba(0, 255, 136, 0.4)',
                    transform: 'translateY(-2px)'
                  },
                  '&:disabled': {
                    background: 'rgba(255, 255, 255, 0.1)',
                    color: '#666666',
                    transform: 'none'
                  }
                }}
                disabled={loading}
              >
                {loading ? '注册中...' : '注册'}
              </Button>
              <Box textAlign="center">
                <Link 
                  component={RouterLink} 
                  to="/login" 
                  variant="body2"
                  sx={{ 
                    color: '#00ff88',
                    textDecoration: 'none',
                    '&:hover': {
                      color: '#00ccff',
                      textDecoration: 'underline'
                    }
                  }}
                >
                  已有账号？立即登录
                </Link>
              </Box>
            </Box>
          </Paper>
        </Box>
      </Container>
    </Box>
  );
};

export default RegisterPage;
