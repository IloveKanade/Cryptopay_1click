import React from 'react';
import {
  Box,
  Container,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  AppBar,
  Toolbar,
  Paper,
} from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import {
  Payment as PaymentIcon,
  Security as SecurityIcon,
  Speed as SpeedIcon,
  Code as CodeIcon,
} from '@mui/icons-material';

const LandingPage: React.FC = () => {
  return (
    <Box sx={{ 
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #16213e 100%)',
      color: 'white'
    }}>
      <AppBar 
        position="static" 
        sx={{ 
          background: 'rgba(0, 0, 0, 0.8)',
          backdropFilter: 'blur(10px)',
          borderBottom: '1px solid rgba(0, 255, 136, 0.3)',
          boxShadow: 'none'
        }}
      >
        <Toolbar>
          <Typography 
            variant="h6" 
            component="div" 
            sx={{ 
              flexGrow: 1,
              color: '#ffffff',
              fontWeight: 'bold'
            }}
          >
            支付链接
          </Typography>
          <Button 
            color="inherit" 
            component={RouterLink} 
            to="/login"
            sx={{ 
              color: '#b0b0b0',
              '&:hover': { color: '#ffffff' }
            }}
          >
            登录
          </Button>
          <Button
            variant="contained"
            component={RouterLink}
            to="/register"
            sx={{ 
              ml: 2,
              background: 'linear-gradient(45deg, #00ff88, #00ccff)',
              color: '#000',
              fontWeight: 'bold',
              '&:hover': {
                background: 'linear-gradient(45deg, #00cc6a, #0099cc)',
              }
            }}
          >
            注册
          </Button>
        </Toolbar>
      </AppBar>

      <Container maxWidth="lg">
        {/* Hero Section */}
        <Box
          sx={{
            textAlign: 'center',
            py: 8,
            mt: 4,
          }}
        >
          <Typography 
            variant="h2" 
            component="h1" 
            gutterBottom
            sx={{ 
              color: '#ffffff',
              fontWeight: 'bold',
              fontSize: { xs: '2.5rem', md: '4rem' }
            }}
          >
            零代码加密货币支付
          </Typography>
          <Typography 
            variant="h5" 
            sx={{ 
              color: '#b0b0b0',
              mb: 4
            }}
          >
            无需编程知识，轻松创建支付链接，接收USDT付款
          </Typography>
          <Button
            variant="contained"
            size="large"
            component={RouterLink}
            to="/register"
            sx={{ 
              mt: 3,
              background: 'linear-gradient(45deg, #00ff88, #00ccff)',
              color: '#000',
              fontWeight: 'bold',
              px: 6,
              py: 2,
              borderRadius: 3,
              textTransform: 'none',
              fontSize: '1.2rem',
              boxShadow: '0 8px 32px rgba(0, 255, 136, 0.3)',
              '&:hover': {
                background: 'linear-gradient(45deg, #00cc6a, #0099cc)',
                boxShadow: '0 12px 40px rgba(0, 255, 136, 0.4)',
                transform: 'translateY(-2px)'
              }
            }}
          >
            立即开始
          </Button>
        </Box>

        {/* Features Section */}
        <Box sx={{ py: 8 }}>
          <Typography 
            variant="h3" 
            component="h2" 
            textAlign="center" 
            gutterBottom
            sx={{ 
              color: '#ffffff',
              mb: 6
            }}
          >
            为什么选择我们？
          </Typography>
          <Grid container spacing={4} sx={{ mt: 4 }}>
            <Grid item xs={12} md={3}>
              <Card sx={{ 
                height: '100%', 
                textAlign: 'center',
                background: 'rgba(26, 26, 46, 0.8)',
                backdropFilter: 'blur(10px)',
                border: '1px solid rgba(0, 255, 136, 0.2)',
                borderRadius: 3,
                transition: 'all 0.3s ease',
                '&:hover': {
                  transform: 'translateY(-8px)',
                  borderColor: 'rgba(0, 255, 136, 0.5)',
                  boxShadow: '0 20px 40px rgba(0, 255, 136, 0.2)'
                }
              }}>
                <CardContent sx={{ p: 3 }}>
                  <CodeIcon sx={{ 
                    fontSize: 60, 
                    color: '#00ff88', 
                    mb: 2,
                    filter: 'drop-shadow(0 0 10px rgba(0, 255, 136, 0.5))'
                  }} />
                  <Typography 
                    variant="h6" 
                    gutterBottom
                    sx={{ 
                      color: '#ffffff',
                      fontWeight: 'bold'
                    }}
                  >
                    零代码集成
                  </Typography>
                  <Typography sx={{ color: '#b0b0b0' }}>
                    无需编写任何代码，只需填写表单即可创建支付链接
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={3}>
              <Card sx={{ 
                height: '100%', 
                textAlign: 'center',
                background: 'rgba(26, 26, 46, 0.8)',
                backdropFilter: 'blur(10px)',
                border: '1px solid rgba(0, 255, 136, 0.2)',
                borderRadius: 3,
                transition: 'all 0.3s ease',
                '&:hover': {
                  transform: 'translateY(-8px)',
                  borderColor: 'rgba(0, 255, 136, 0.5)',
                  boxShadow: '0 20px 40px rgba(0, 255, 136, 0.2)'
                }
              }}>
                <CardContent sx={{ p: 3 }}>
                  <SecurityIcon sx={{ 
                    fontSize: 60, 
                    color: '#00ff88', 
                    mb: 2,
                    filter: 'drop-shadow(0 0 10px rgba(0, 255, 136, 0.5))'
                  }} />
                  <Typography 
                    variant="h6" 
                    gutterBottom
                    sx={{ 
                      color: '#ffffff',
                      fontWeight: 'bold'
                    }}
                  >
                    安全可靠
                  </Typography>
                  <Typography sx={{ color: '#b0b0b0' }}>
                    基于区块链技术，交易安全透明，资金直接到账
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={3}>
              <Card sx={{ 
                height: '100%', 
                textAlign: 'center',
                background: 'rgba(26, 26, 46, 0.8)',
                backdropFilter: 'blur(10px)',
                border: '1px solid rgba(0, 255, 136, 0.2)',
                borderRadius: 3,
                transition: 'all 0.3s ease',
                '&:hover': {
                  transform: 'translateY(-8px)',
                  borderColor: 'rgba(0, 255, 136, 0.5)',
                  boxShadow: '0 20px 40px rgba(0, 255, 136, 0.2)'
                }
              }}>
                <CardContent sx={{ p: 3 }}>
                  <SpeedIcon sx={{ 
                    fontSize: 60, 
                    color: '#00ff88', 
                    mb: 2,
                    filter: 'drop-shadow(0 0 10px rgba(0, 255, 136, 0.5))'
                  }} />
                  <Typography 
                    variant="h6" 
                    gutterBottom
                    sx={{ 
                      color: '#ffffff',
                      fontWeight: 'bold'
                    }}
                  >
                    快速支付
                  </Typography>
                  <Typography sx={{ color: '#b0b0b0' }}>
                    支持USDT支付，秒级确认，全球通用
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={3}>
              <Card sx={{ 
                height: '100%', 
                textAlign: 'center',
                background: 'rgba(26, 26, 46, 0.8)',
                backdropFilter: 'blur(10px)',
                border: '1px solid rgba(0, 255, 136, 0.2)',
                borderRadius: 3,
                transition: 'all 0.3s ease',
                '&:hover': {
                  transform: 'translateY(-8px)',
                  borderColor: 'rgba(0, 255, 136, 0.5)',
                  boxShadow: '0 20px 40px rgba(0, 255, 136, 0.2)'
                }
              }}>
                <CardContent sx={{ p: 3 }}>
                  <PaymentIcon sx={{ 
                    fontSize: 60, 
                    color: '#00ff88', 
                    mb: 2,
                    filter: 'drop-shadow(0 0 10px rgba(0, 255, 136, 0.5))'
                  }} />
                  <Typography 
                    variant="h6" 
                    gutterBottom
                    sx={{ 
                      color: '#ffffff',
                      fontWeight: 'bold'
                    }}
                  >
                    简单易用
                  </Typography>
                  <Typography sx={{ color: '#b0b0b0' }}>
                    直观的用户界面，一键生成支付链接
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Box>

        {/* How it works */}
        <Box sx={{ py: 8 }}>
          <Typography 
            variant="h3" 
            component="h2" 
            textAlign="center" 
            gutterBottom
            sx={{ 
              color: '#ffffff',
              mb: 6
            }}
          >
            如何使用？
          </Typography>
          <Grid container spacing={4} sx={{ mt: 4 }}>
            <Grid item xs={12} md={4}>
              <Box textAlign="center">
                <Paper sx={{ 
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: 80,
                  height: 80,
                  borderRadius: '50%',
                  background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                  color: '#000',
                  fontWeight: 'bold',
                  fontSize: '2rem',
                  mb: 3,
                  boxShadow: '0 8px 32px rgba(0, 255, 136, 0.3)'
                }}>
                  1
                </Paper>
                <Typography 
                  variant="h6" 
                  gutterBottom
                  sx={{ 
                    color: '#ffffff',
                    fontWeight: 'bold'
                  }}
                >
                  创建支付链接
                </Typography>
                <Typography sx={{ color: '#b0b0b0' }}>
                  登录后台，填写金额和钱包地址，一键生成支付链接
                </Typography>
              </Box>
            </Grid>
            <Grid item xs={12} md={4}>
              <Box textAlign="center">
                <Paper sx={{ 
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: 80,
                  height: 80,
                  borderRadius: '50%',
                  background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                  color: '#000',
                  fontWeight: 'bold',
                  fontSize: '2rem',
                  mb: 3,
                  boxShadow: '0 8px 32px rgba(0, 255, 136, 0.3)'
                }}>
                  2
                </Paper>
                <Typography 
                  variant="h6" 
                  gutterBottom
                  sx={{ 
                    color: '#ffffff',
                    fontWeight: 'bold'
                  }}
                >
                  分享链接
                </Typography>
                <Typography sx={{ color: '#b0b0b0' }}>
                  将生成的链接分享给客户，或嵌入到您的网站中
                </Typography>
              </Box>
            </Grid>
            <Grid item xs={12} md={4}>
              <Box textAlign="center">
                <Paper sx={{ 
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: 80,
                  height: 80,
                  borderRadius: '50%',
                  background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                  color: '#000',
                  fontWeight: 'bold',
                  fontSize: '2rem',
                  mb: 3,
                  boxShadow: '0 8px 32px rgba(0, 255, 136, 0.3)'
                }}>
                  3
                </Paper>
                <Typography 
                  variant="h6" 
                  gutterBottom
                  sx={{ 
                    color: '#ffffff',
                    fontWeight: 'bold'
                  }}
                >
                  接收付款
                </Typography>
                <Typography sx={{ color: '#b0b0b0' }}>
                  客户通过链接完成支付，资金直接到您的钱包
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </Box>

        {/* CTA Section */}
        <Box
          sx={{
            textAlign: 'center',
            py: 8,
            background: 'rgba(26, 26, 46, 0.8)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(0, 255, 136, 0.3)',
            borderRadius: 3,
            mt: 4,
          }}
        >
          <Typography 
            variant="h4" 
            gutterBottom
            sx={{ 
              color: '#ffffff'
            }}
          >
            准备开始了吗？
          </Typography>
          <Typography 
            variant="body1" 
            sx={{ 
              color: '#b0b0b0',
              mb: 4
            }}
          >
            加入我们，体验最简单的加密货币支付解决方案
          </Typography>
          <Button
            variant="contained"
            size="large"
            component={RouterLink}
            to="/register"
            sx={{
              background: 'linear-gradient(45deg, #00ff88, #00ccff)',
              color: '#000',
              fontWeight: 'bold',
              px: 6,
              py: 2,
              borderRadius: 3,
              textTransform: 'none',
              fontSize: '1.2rem',
              boxShadow: '0 8px 32px rgba(0, 255, 136, 0.3)',
              '&:hover': {
                background: 'linear-gradient(45deg, #00cc6a, #0099cc)',
                boxShadow: '0 12px 40px rgba(0, 255, 136, 0.4)',
                transform: 'translateY(-2px)'
              }
            }}
          >
            免费注册
          </Button>
        </Box>
      </Container>
    </Box>
  );
};

export default LandingPage;
