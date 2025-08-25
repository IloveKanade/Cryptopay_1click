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
  Divider,
} from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // 使用AuthContext的login函数
      const userData = await login(email, password);
      
      // 根据用户角色决定跳转路径
      if (userData.role === 'admin') {
        navigate('/admin');
      } else {
        navigate('/dashboard');
      }
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
              登录
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
                id="email"
                label="邮箱地址"
                name="email"
                autoComplete="email"
                autoFocus
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                sx={{
                  '& .MuiOutlinedInput-root': {
                    '& fieldset': {
                      borderColor: 'rgba(0, 255, 136, 0.3)',
                    },
                    '&:hover fieldset': {
                      borderColor: 'rgba(0, 255, 136, 0.5)',
                    },
                    '&.Mui-focused fieldset': {
                      borderColor: 'rgba(0, 255, 136, 0.8)',
                    },
                  },
                  '& .MuiInputLabel-root': {
                    color: 'rgba(255, 255, 255, 0.7)',
                  },
                  '& .MuiInputBase-input': {
                    color: '#ffffff',
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
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                sx={{
                  '& .MuiOutlinedInput-root': {
                    '& fieldset': {
                      borderColor: 'rgba(0, 255, 136, 0.3)',
                    },
                    '&:hover fieldset': {
                      borderColor: 'rgba(0, 255, 136, 0.5)',
                    },
                    '&.Mui-focused fieldset': {
                      borderColor: 'rgba(0, 255, 136, 0.8)',
                    },
                  },
                  '& .MuiInputLabel-root': {
                    color: 'rgba(255, 255, 255, 0.7)',
                  },
                  '& .MuiInputBase-input': {
                    color: '#ffffff',
                  },
                }}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                disabled={loading}
                sx={{
                  mt: 3,
                  mb: 2,
                  py: 1.5,
                  background: 'linear-gradient(45deg, #00ff88 30%, #00cc6a 90%)',
                  color: '#000',
                  fontWeight: 'bold',
                  fontSize: '1.1rem',
                  '&:hover': {
                    background: 'linear-gradient(45deg, #00cc6a 30%, #00ff88 90%)',
                    transform: 'translateY(-2px)',
                    boxShadow: '0 8px 25px rgba(0, 255, 136, 0.3)',
                  },
                  '&:disabled': {
                    background: 'rgba(0, 255, 136, 0.3)',
                    color: 'rgba(255, 255, 255, 0.5)',
                  },
                }}
              >
                {loading ? '登录中...' : '登录'}
              </Button>
              
              <Divider sx={{ my: 2, borderColor: 'rgba(255, 255, 255, 0.2)' }}>
                <Typography variant="body2" sx={{ color: 'rgba(255, 255, 255, 0.6)' }}>
                  或
                </Typography>
              </Divider>
              
              <Box sx={{ textAlign: 'center' }}>
                <Link
                  component={RouterLink}
                  to="/register"
                  variant="body2"
                  sx={{
                    color: 'rgba(0, 255, 136, 0.8)',
                    textDecoration: 'none',
                    '&:hover': {
                      color: 'rgba(0, 255, 136, 1)',
                      textDecoration: 'underline',
                    },
                  }}
                >
                  还没有账户？立即注册
                </Link>
              </Box>
            </Box>
          </Paper>
        </Box>
      </Container>
    </Box>
  );
};

export default LoginPage;
