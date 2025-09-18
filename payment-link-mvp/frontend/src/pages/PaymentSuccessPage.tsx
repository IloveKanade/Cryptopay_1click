import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  Card,
  CardContent,
  Alert,
  CircularProgress,
  Divider,
  Grid,
} from '@mui/material';
import {
  CheckCircle as CheckCircleIcon,
  Home as HomeIcon,
  Receipt as ReceiptIcon,
} from '@mui/icons-material';
import { api } from '../services/api';

interface PaymentOrder {
  id: number;
  order_id: string;
  trade_id: string;
  amount: number;
  actual_amount: number;
  token: string;
  status: string;
  payment_url: string;
  expiration_time: string;
  block_transaction_id: string;
  created_at: string;
  updated_at: string;
  payment_link: {
    id: number;
    title: string;
    description: string;
    currency: string;
  };
}

const PaymentSuccessPage: React.FC = () => {
  const { orderId } = useParams<{ orderId: string }>();
  const navigate = useNavigate();
  const [paymentOrder, setPaymentOrder] = useState<PaymentOrder | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (orderId) {
      fetchPaymentOrder();
    }
  }, [orderId]);

  const fetchPaymentOrder = async () => {
    try {
      const response = await api.get(`/api/payment/orders/${orderId}`);
      setPaymentOrder(response.data.data);
    } catch (err: any) {
      setError(err.response?.data?.error || '获取订单信息失败');
    } finally {
      setLoading(false);
    }
  };

  const handleGoHome = () => {
    navigate('/');
  };

  const handleViewReceipt = () => {
    // 这里可以添加查看收据的功能
    console.log('查看收据');
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

  if (error || !paymentOrder) {
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
            {error || '订单信息不存在'}
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
          {/* 成功图标和标题 */}
          <Box textAlign="center" sx={{ mb: 4 }}>
            <CheckCircleIcon sx={{ 
              fontSize: 80, 
              color: '#00ff88', 
              mb: 2,
              filter: 'drop-shadow(0 0 20px rgba(0, 255, 136, 0.5))'
            }} />
            <Typography 
              variant="h3" 
              component="h1" 
              gutterBottom
              sx={{ 
                background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                backgroundClip: 'text',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                fontWeight: 'bold'
              }}
            >
              支付成功！
            </Typography>
            <Typography 
              variant="h6" 
              sx={{ 
                color: '#b0b0b0'
              }}
            >
              您的USDT支付已完成
            </Typography>
          </Box>

          <Divider sx={{ 
            my: 4,
            borderColor: 'rgba(0, 255, 136, 0.2)'
          }} />

          {/* 订单详情 */}
          <Card sx={{ 
            mb: 4,
            background: 'rgba(0, 0, 0, 0.3)',
            border: '1px solid rgba(0, 255, 136, 0.2)',
            borderRadius: 3
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
                订单详情
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" sx={{ color: '#b0b0b0' }}>
                    订单号
                  </Typography>
                  <Typography variant="body1" sx={{ 
                    fontFamily: 'monospace',
                    color: '#ffffff'
                  }}>
                    {paymentOrder.order_id}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" sx={{ color: '#b0b0b0' }}>
                    交易号
                  </Typography>
                  <Typography variant="body1" sx={{ 
                    fontFamily: 'monospace',
                    color: '#ffffff'
                  }}>
                    {paymentOrder.trade_id}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" sx={{ color: '#b0b0b0' }}>
                    支付金额
                  </Typography>
                  <Typography 
                    variant="h6"
                    sx={{ 
                      background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                      backgroundClip: 'text',
                      WebkitBackgroundClip: 'text',
                      WebkitTextFillColor: 'transparent',
                      fontWeight: 'bold'
                    }}
                  >
                    {paymentOrder.amount} CNY
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" sx={{ color: '#b0b0b0' }}>
                    实际USDT
                  </Typography>
                  <Typography 
                    variant="h6"
                    sx={{ 
                      background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                      backgroundClip: 'text',
                      WebkitBackgroundClip: 'text',
                      WebkitTextFillColor: 'transparent',
                      fontWeight: 'bold'
                    }}
                  >
                    {paymentOrder.actual_amount} USDT
                  </Typography>
                </Grid>
                {paymentOrder.block_transaction_id && (
                  <Grid item xs={12}>
                    <Typography variant="body2" sx={{ color: '#b0b0b0' }}>
                      区块链交易ID
                    </Typography>
                    <Typography variant="body1" sx={{ 
                      fontFamily: 'monospace', 
                      wordBreak: 'break-all',
                      color: '#ffffff'
                    }}>
                      {paymentOrder.block_transaction_id}
                    </Typography>
                  </Grid>
                )}
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" sx={{ color: '#b0b0b0' }}>
                    支付时间
                  </Typography>
                  <Typography variant="body1" sx={{ color: '#ffffff' }}>
                    {new Date(paymentOrder.updated_at).toLocaleString('zh-CN')}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" sx={{ color: '#b0b0b0' }}>
                    支付状态
                  </Typography>
                  <Typography 
                    variant="body1" 
                    sx={{ 
                      color: '#00ff88',
                      fontWeight: 'bold'
                    }}
                  >
                    {paymentOrder.status === 'paid' ? '已支付' : paymentOrder.status}
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>

          {/* 支付链接信息 */}
          {paymentOrder.payment_link && (
            <Card sx={{ 
              mb: 4,
              background: 'rgba(0, 0, 0, 0.3)',
              border: '1px solid rgba(0, 255, 136, 0.2)',
              borderRadius: 3
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
                  支付链接信息
                </Typography>
                <Typography variant="body1" gutterBottom sx={{ color: '#ffffff' }}>
                  <strong style={{ color: '#00ff88' }}>标题：</strong>{paymentOrder.payment_link.title}
                </Typography>
                {paymentOrder.payment_link.description && (
                  <Typography variant="body1" gutterBottom sx={{ color: '#ffffff' }}>
                    <strong style={{ color: '#00ff88' }}>描述：</strong>{paymentOrder.payment_link.description}
                  </Typography>
                )}
                <Typography variant="body1" sx={{ color: '#ffffff' }}>
                  <strong style={{ color: '#00ff88' }}>货币：</strong>{paymentOrder.payment_link.currency}
                </Typography>
              </CardContent>
            </Card>
          )}

          {/* 操作按钮 */}
          <Box textAlign="center" sx={{ mt: 4 }}>
            <Button
              variant="contained"
              startIcon={<HomeIcon />}
              onClick={handleGoHome}
              sx={{ 
                mr: 2,
                background: 'linear-gradient(45deg, #00ff88, #00ccff)',
                color: '#000',
                fontWeight: 'bold',
                px: 4,
                py: 1.5,
                borderRadius: 3,
                textTransform: 'none',
                fontSize: '1.1rem',
                boxShadow: '0 8px 32px rgba(0, 255, 136, 0.3)',
                '&:hover': {
                  background: 'linear-gradient(45deg, #00cc6a, #0099cc)',
                  boxShadow: '0 12px 40px rgba(0, 255, 136, 0.4)',
                  transform: 'translateY(-2px)'
                }
              }}
            >
              返回首页
            </Button>
            <Button
              variant="outlined"
              startIcon={<ReceiptIcon />}
              onClick={handleViewReceipt}
              sx={{ 
                color: '#00ff88',
                borderColor: 'rgba(0, 255, 136, 0.3)',
                px: 4,
                py: 1.5,
                borderRadius: 3,
                textTransform: 'none',
                fontSize: '1.1rem',
                '&:hover': {
                  borderColor: 'rgba(0, 255, 136, 0.5)',
                  backgroundColor: 'rgba(0, 255, 136, 0.1)',
                  color: '#ffffff'
                }
              }}
            >
              查看收据
            </Button>
          </Box>

          {/* 温馨提示 */}
          <Box sx={{ 
            background: 'rgba(0, 255, 136, 0.1)', 
            border: '1px solid rgba(0, 255, 136, 0.3)',
            p: 3, 
            borderRadius: 3, 
            mt: 4 
          }}>
            <Typography 
              variant="h6" 
              gutterBottom
              sx={{ 
                color: '#00ff88',
                fontWeight: 'bold'
              }}
            >
              温馨提示
            </Typography>
            <Typography variant="body2" paragraph sx={{ color: '#b0b0b0' }}>
              • 您的支付已成功完成，资金已转入收款方钱包
            </Typography>
            <Typography variant="body2" paragraph sx={{ color: '#b0b0b0' }}>
              • 请保存好订单号，以便后续查询
            </Typography>
            <Typography variant="body2" paragraph sx={{ color: '#b0b0b0' }}>
              • 如有任何问题，请联系客服
            </Typography>
          </Box>
        </Paper>
      </Container>
    </Box>
  );
};

export default PaymentSuccessPage;
