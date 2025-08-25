import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  Alert,
  CircularProgress,
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import {
  Payment as PaymentIcon,
  ContentCopy as CopyIcon,
  QrCode as QrCodeIcon,
  CheckCircle as CheckCircleIcon,
  Error as ErrorIcon,
  Timer as TimerIcon,
} from '@mui/icons-material';
import { api } from '../services/api';

interface PaymentLink {
  id: number;
  title: string;
  amount: number;
  currency: string;
  wallet_address: string;
  link_id: string;
  description: string;
  status: string;
  created_at: string;
}

interface PaymentOrder {
  order_id: string;
  trade_id: string;
  amount: number;
  actual_amount: number;
  payment_url: string;
  expires_at: string;
  status: string;
}

const PaymentPage: React.FC = () => {
  const { linkId } = useParams<{ linkId: string }>();
  const navigate = useNavigate();
  const [paymentLink, setPaymentLink] = useState<PaymentLink | null>(null);
  const [paymentOrder, setPaymentOrder] = useState<PaymentOrder | null>(null);
  const [loading, setLoading] = useState(true);
  const [creatingOrder, setCreatingOrder] = useState(false);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);
  const [showPaymentDialog, setShowPaymentDialog] = useState(false);
  const [orderStatus, setOrderStatus] = useState<string>('pending');
  const [statusCheckInterval, setStatusCheckInterval] = useState<NodeJS.Timeout | null>(null);

  useEffect(() => {
    if (linkId) {
      fetchPaymentLink();
    }
  }, [linkId]);

  useEffect(() => {
    // 清理定时器
    return () => {
      if (statusCheckInterval) {
        clearInterval(statusCheckInterval);
      }
    };
  }, [statusCheckInterval]);

  const fetchPaymentLink = async () => {
    try {
      const response = await api.get(`/api/pay/${linkId}`);
      setPaymentLink(response.data.payment_link);
    } catch (err: any) {
      setError(err.response?.data?.error || '获取支付信息失败');
    } finally {
      setLoading(false);
    }
  };

  const createPaymentOrder = async () => {
    if (!paymentLink) return;

    setCreatingOrder(true);
    setError('');

    try {
      const response = await api.post('/api/payment/create', {
        payment_link_id: paymentLink.id,
        amount: paymentLink.amount,
        client_ip: '127.0.0.1', // 在实际环境中应该获取真实IP
      });

      const orderData = response.data.data;
      setPaymentOrder(orderData);
      setShowPaymentDialog(true);
      
      // 开始轮询订单状态
      startStatusCheck(orderData.trade_id);
    } catch (err: any) {
      setError(err.response?.data?.error || '创建支付订单失败');
    } finally {
      setCreatingOrder(false);
    }
  };

  const startStatusCheck = (tradeId: string) => {
    const interval = setInterval(async () => {
      try {
        const response = await api.get(`/api/payment/status/${tradeId}`);
        const status = response.data.data.status;
        setOrderStatus(status);

        if (status === 'paid') {
          clearInterval(interval);
          setStatusCheckInterval(null);
          // 支付成功，跳转到成功页面
          setTimeout(() => {
            navigate(`/payment/success/${paymentOrder?.order_id}`);
          }, 2000);
        } else if (status === 'expired') {
          clearInterval(interval);
          setStatusCheckInterval(null);
          setError('支付订单已过期');
        }
      } catch (err) {
        console.error('查询订单状态失败:', err);
      }
    }, 3000); // 每3秒查询一次

    setStatusCheckInterval(interval);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const openPaymentWindow = (paymentUrl: string) => {
    const width = 500;
    const height = 600;
    const left = window.screenX + (window.outerWidth - width) / 2;
    const top = window.screenY + (window.outerHeight - height) / 2;

    window.open(
      paymentUrl,
      'payment',
      `width=${width},height=${height},left=${left},top=${top},scrollbars=yes,resizable=yes`
    );
  };

  const getStatusIcon = () => {
    switch (orderStatus) {
      case 'paid':
        return <CheckCircleIcon sx={{ color: '#00ff88', fontSize: 40 }} />;
      case 'expired':
        return <ErrorIcon sx={{ color: '#ff4444', fontSize: 40 }} />;
      case 'pending':
        return <TimerIcon sx={{ color: '#ffaa00', fontSize: 40 }} />;
      default:
        return <TimerIcon sx={{ color: '#ffaa00', fontSize: 40 }} />;
    }
  };

  const getStatusText = () => {
    switch (orderStatus) {
      case 'paid':
        return '支付成功';
      case 'expired':
        return '订单已过期';
      case 'pending':
        return '等待支付';
      default:
        return '等待支付';
    }
  };

  if (loading) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
          background: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #16213e 100%)',
        }}
      >
        <CircularProgress sx={{ color: '#00ff88' }} />
      </Box>
    );
  }

  if (error || !paymentLink) {
    return (
      <Box sx={{ 
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #16213e 100%)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        py: 4
      }}>
        <Container maxWidth="md">
          <Alert 
            severity="error"
            sx={{ 
              background: 'rgba(255, 68, 68, 0.1)',
              border: '1px solid rgba(255, 68, 68, 0.3)',
              color: '#ff4444'
            }}
          >
            {error || '支付链接不存在或已失效'}
          </Alert>
        </Container>
      </Box>
    );
  }

  return (
    <Box sx={{ 
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #16213e 100%)',
      color: 'white',
      py: 4
    }}>
      <Container maxWidth="md">
        <Paper 
          elevation={0}
          sx={{ 
            p: 4,
            background: 'rgba(26, 26, 46, 0.9)',
            backdropFilter: 'blur(20px)',
            border: '1px solid rgba(0, 255, 136, 0.3)',
            borderRadius: 3,
            boxShadow: '0 20px 40px rgba(0, 0, 0, 0.5)'
          }}
        >
          <Box textAlign="center" sx={{ mb: 4 }}>
            <PaymentIcon sx={{ 
              fontSize: 60, 
              color: '#00ff88', 
              mb: 2,
              filter: 'drop-shadow(0 0 10px rgba(0, 255, 136, 0.5))'
            }} />
            <Typography 
              variant="h4" 
              component="h1" 
              gutterBottom
              sx={{ 
                color: '#ffffff',
                fontWeight: 'bold'
              }}
            >
              {paymentLink.title}
            </Typography>
            {paymentLink.description && (
              <Typography 
                variant="body1" 
                sx={{ 
                  color: '#b0b0b0',
                  mb: 2,
                  fontStyle: 'italic'
                }}
              >
                {paymentLink.description}
              </Typography>
            )}
          </Box>

          <Grid container spacing={4}>
            {/* 支付金额 */}
            <Grid item xs={12} md={6}>
              <Card sx={{
                background: 'rgba(0, 0, 0, 0.3)',
                border: '1px solid rgba(0, 255, 136, 0.2)',
                borderRadius: 3,
                transition: 'all 0.3s ease',
                '&:hover': {
                  borderColor: 'rgba(0, 255, 136, 0.5)',
                  boxShadow: '0 10px 30px rgba(0, 255, 136, 0.2)'
                }
              }}>
                <CardContent sx={{ p: 3 }}>
                  <Typography 
                    variant="h6" 
                    gutterBottom
                    sx={{ 
                      color: '#ffffff',
                      fontWeight: 'bold'
                    }}
                  >
                    支付金额
                  </Typography>
                  <Typography 
                    variant="h3" 
                    sx={{ 
                      background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                      backgroundClip: 'text',
                      WebkitBackgroundClip: 'text',
                      WebkitTextFillColor: 'transparent',
                      fontWeight: 'bold',
                      mb: 2
                    }}
                  >
                    {paymentLink.amount} {paymentLink.currency}
                  </Typography>
                  <Typography 
                    variant="body2" 
                    sx={{ color: '#b0b0b0' }}
                  >
                    请使用USDT进行支付
                  </Typography>
                </CardContent>
              </Card>
            </Grid>

            {/* 钱包地址 */}
            <Grid item xs={12} md={6}>
              <Card sx={{
                background: 'rgba(0, 0, 0, 0.3)',
                border: '1px solid rgba(0, 255, 136, 0.2)',
                borderRadius: 3,
                transition: 'all 0.3s ease',
                '&:hover': {
                  borderColor: 'rgba(0, 255, 136, 0.5)',
                  boxShadow: '0 10px 30px rgba(0, 255, 136, 0.2)'
                }
              }}>
                <CardContent sx={{ p: 3 }}>
                  <Typography 
                    variant="h6" 
                    gutterBottom
                    sx={{ 
                      color: '#ffffff',
                      fontWeight: 'bold'
                    }}
                  >
                    收款钱包地址
                  </Typography>
                  <Box
                    sx={{
                      p: 2,
                      background: 'rgba(0, 0, 0, 0.5)',
                      border: '1px solid rgba(0, 255, 136, 0.1)',
                      borderRadius: 2,
                      wordBreak: 'break-all',
                      fontFamily: 'monospace',
                      fontSize: '0.875rem',
                      color: '#b0b0b0',
                      mb: 2
                    }}
                  >
                    {paymentLink.wallet_address}
                  </Box>
                  <Button
                    variant="outlined"
                    startIcon={<CopyIcon />}
                    onClick={() => copyToClipboard(paymentLink.wallet_address)}
                    sx={{ 
                      mt: 2,
                      color: '#00ff88',
                      borderColor: 'rgba(0, 255, 136, 0.3)',
                      '&:hover': {
                        borderColor: 'rgba(0, 255, 136, 0.5)',
                        backgroundColor: 'rgba(0, 255, 136, 0.1)'
                      }
                    }}
                    fullWidth
                  >
                    {copied ? '已复制' : '复制地址'}
                  </Button>
                </CardContent>
              </Card>
            </Grid>
          </Grid>

          <Divider sx={{ 
            my: 4,
            borderColor: 'rgba(0, 255, 136, 0.2)'
          }} />

          {/* 支付按钮 */}
          <Box textAlign="center" sx={{ mb: 4 }}>
            <Button
              variant="contained"
              size="large"
              startIcon={<PaymentIcon />}
              onClick={createPaymentOrder}
              disabled={creatingOrder}
              sx={{ 
                px: 6, 
                py: 2, 
                fontSize: '1.2rem',
                background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                color: '#000',
                fontWeight: 'bold',
                borderRadius: 3,
                textTransform: 'none',
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
            >
              {creatingOrder ? '创建订单中...' : '立即支付'}
            </Button>
          </Box>

          {/* 支付说明 */}
          <Box sx={{ 
            background: 'rgba(0, 0, 0, 0.3)', 
            border: '1px solid rgba(0, 255, 136, 0.2)',
            p: 3, 
            borderRadius: 3 
          }}>
            <Typography 
              variant="h6" 
              gutterBottom
              sx={{ 
                color: '#ffffff',
                fontWeight: 'bold'
              }}
            >
              支付说明
            </Typography>
            <Typography variant="body2" paragraph sx={{ color: '#b0b0b0' }}>
              1. 点击"立即支付"按钮创建支付订单
            </Typography>
            <Typography variant="body2" paragraph sx={{ color: '#b0b0b0' }}>
              2. 系统将跳转到Epusdt收银台页面
            </Typography>
            <Typography variant="body2" paragraph sx={{ color: '#b0b0b0' }}>
              3. 在收银台页面完成USDT支付
            </Typography>
            <Typography 
              variant="body2" 
              sx={{ 
                color: '#ffaa00',
                fontWeight: 'bold'
              }}
            >
              注意：请确保支付金额与显示金额一致，支付完成后请勿关闭页面
            </Typography>
          </Box>
        </Paper>

        {/* 支付状态对话框 */}
        <Dialog
          open={showPaymentDialog}
          onClose={() => setShowPaymentDialog(false)}
          maxWidth="sm"
          fullWidth
          PaperProps={{
            sx: {
              background: 'rgba(26, 26, 46, 0.95)',
              backdropFilter: 'blur(20px)',
              border: '1px solid rgba(0, 255, 136, 0.3)',
              borderRadius: 3
            }
          }}
        >
          <DialogTitle sx={{ 
            color: '#ffffff',
            borderBottom: '1px solid rgba(0, 255, 136, 0.2)'
          }}>
            <Box display="flex" alignItems="center" gap={2}>
              {getStatusIcon()}
              <Typography variant="h6">支付状态</Typography>
            </Box>
          </DialogTitle>
          <DialogContent sx={{ pt: 3 }}>
            <Box textAlign="center" sx={{ py: 2 }}>
              <Typography 
                variant="h5" 
                gutterBottom
                sx={{ 
                  color: '#ffffff',
                  fontWeight: 'bold'
                }}
              >
                {getStatusText()}
              </Typography>
              {paymentOrder && (
                <>
                  <Typography variant="body1" sx={{ color: '#b0b0b0' }} gutterBottom>
                    订单号: {paymentOrder.order_id}
                  </Typography>
                  <Typography variant="body1" sx={{ color: '#b0b0b0' }} gutterBottom>
                    交易号: {paymentOrder.trade_id}
                  </Typography>
                  <Typography 
                    variant="h6" 
                    sx={{ 
                      background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                      backgroundClip: 'text',
                      WebkitBackgroundClip: 'text',
                      WebkitTextFillColor: 'transparent',
                      fontWeight: 'bold',
                      mb: 1
                    }}
                  >
                    支付金额: {paymentOrder.amount} CNY
                  </Typography>
                  <Typography variant="body2" sx={{ color: '#b0b0b0' }} gutterBottom>
                    实际USDT: {paymentOrder.actual_amount} USDT
                  </Typography>
                </>
              )}
              
              {orderStatus === 'pending' && (
                <Box sx={{ mt: 3 }}>
                  <Button
                    variant="contained"
                    onClick={() => paymentOrder && openPaymentWindow(paymentOrder.payment_url)}
                    sx={{ 
                      mr: 2,
                      background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                      color: '#000',
                      fontWeight: 'bold',
                      '&:hover': {
                        background: 'linear-gradient(45deg, #00cc6a, #0099cc)',
                      }
                    }}
                  >
                    打开收银台
                  </Button>
                  <Button
                    variant="outlined"
                    onClick={() => setShowPaymentDialog(false)}
                    sx={{ 
                      color: '#b0b0b0',
                      borderColor: 'rgba(0, 255, 136, 0.3)',
                      '&:hover': {
                        borderColor: 'rgba(0, 255, 136, 0.5)',
                        color: '#ffffff'
                      }
                    }}
                  >
                    关闭
                  </Button>
                </Box>
              )}
            </Box>
          </DialogContent>
        </Dialog>
      </Container>
    </Box>
  );
};

export default PaymentPage;
