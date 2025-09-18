import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Box,
  Button,
  Card,
  CardContent,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Chip,
  Alert,
  CircularProgress,
  Switch,
  FormControlLabel,
  Divider,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Visibility as ViewIcon,
  AccountBalanceWallet as WalletIcon,
  Settings as SettingsIcon,
  Receipt as ReceiptIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { api } from '../services/api';
import AdminNavigation from '../components/AdminNavigation';

interface MiddlemanWallet {
  id: number;
  wallet_address: string;
  name: string;
  status: 'active' | 'inactive';
  balance: number;
  created_at: string;
  updated_at: string;
}

interface TransferRecord {
  id: number;
  payment_order_id: number;
  payment_order: {
    order_id: string;
    trade_id: string;
  };
  middleman_wallet_id: number;
  middleman_wallet: {
    name: string;
    wallet_address: string;
  };
  from_address: string;
  to_address: string;
  amount: number;
  fixed_fee_amount: number;
  network_fee_amount: number;
  total_fee_amount: number;
  net_amount: number;
  network_fees: string;
  status: 'pending' | 'success' | 'failed';
  transaction_hash: string;
  error_message: string;
  created_at: string;
  updated_at: string;
}

interface FeeConfig {
  id: number;
  fee_rate: number;
  min_fee: number;
  max_fee: number;
  network_fee_min: number;
  network_fee_max: number;
  created_at: string;
  updated_at: string;
}

const AdminMiddlemanPage: React.FC = () => {
  const [wallets, setWallets] = useState<MiddlemanWallet[]>([]);
  const [transfers, setTransfers] = useState<TransferRecord[]>([]);
  const [feeConfig, setFeeConfig] = useState<FeeConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'wallets' | 'transfers' | 'fees'>('wallets');
  
  // 对话框状态
  const [walletDialogOpen, setWalletDialogOpen] = useState(false);
  const [feeDialogOpen, setFeeDialogOpen] = useState(false);
  const [newWallet, setNewWallet] = useState({ name: '', wallet_address: '', private_key: '' });
  const [newFeeConfig, setNewFeeConfig] = useState({ 
    fee_rate: 1, 
    min_fee: 0.1, 
    max_fee: 10,
    network_fee_min: 0.001,
    network_fee_max: 0.01
  });
  
  const navigate = useNavigate();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // 并行获取数据
      const [walletsRes, transfersRes, feeConfigRes] = await Promise.all([
        api.get('/api/admin/middleman-wallets'),
        api.get('/api/admin/transfer-records'),
        api.get('/api/admin/fee-config')
      ]);
      
      setWallets(walletsRes.data.data || []);
      setTransfers(transfersRes.data.data || []);
      setFeeConfig(feeConfigRes.data.data);
    } catch (err: any) {
      setError(err.response?.data?.error || '获取数据失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAddWallet = async () => {
    try {
      // 直接使用正确的字段名
      const walletData = {
        name: newWallet.name,
        wallet_address: newWallet.wallet_address,
        private_key: newWallet.private_key
      };
      await api.post('/api/admin/middleman-wallets', walletData);
      setWalletDialogOpen(false);
      setNewWallet({ name: '', wallet_address: '', private_key: '' });
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || '添加钱包失败');
    }
  };

  const handleUpdateWalletStatus = async (id: number, status: 'active' | 'inactive') => {
    try {
      await api.put(`/api/admin/middleman-wallets/${id}/status`, { status });
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || '更新钱包状态失败');
    }
  };

  const handleDeleteWallet = async (id: number) => {
    if (!window.confirm('确定要删除这个钱包吗？')) return;
    
    try {
      await api.delete(`/api/admin/middleman-wallets/${id}`);
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || '删除钱包失败');
    }
  };

  const handleUpdateFeeConfig = async () => {
    try {
      if (feeConfig) {
        await api.put('/api/admin/fee-config', newFeeConfig);
      } else {
        await api.post('/api/admin/fee-config', newFeeConfig);
      }
      setFeeDialogOpen(false);
      fetchData();
    } catch (err: any) {
      setError(err.response?.data?.error || '更新费率配置失败');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
      case 'success':
      case 'completed':
        return 'success';
      case 'pending':
        return 'warning';
      case 'inactive':
      case 'failed':
        return 'error';
      default:
        return 'default';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'active':
        return '启用';
      case 'inactive':
        return '禁用';
      case 'pending':
        return '处理中';
      case 'success':
      case 'completed':
        return '已完成';
      case 'failed':
        return '失败';
      default:
        return status;
    }
  };

  if (loading) {
    return (
      <>
        <AdminNavigation />
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
            <CircularProgress />
          </Box>
        </Container>
      </>
    );
  }

  return (
    <>
      <AdminNavigation />
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        {/* 页面标题 */}
        <Box display="flex" alignItems="center" mb={4}>
          <WalletIcon sx={{ mr: 2, fontSize: 32, color: 'primary.main' }} />
          <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold' }}>
            中间人钱包管理
          </Typography>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        {/* 标签页导航 */}
        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button
              variant={activeTab === 'wallets' ? 'contained' : 'outlined'}
              onClick={() => setActiveTab('wallets')}
              startIcon={<WalletIcon />}
            >
              钱包管理
            </Button>
            <Button
              variant={activeTab === 'transfers' ? 'contained' : 'outlined'}
              onClick={() => setActiveTab('transfers')}
              startIcon={<ReceiptIcon />}
            >
              转账记录
            </Button>
            <Button
              variant={activeTab === 'fees' ? 'contained' : 'outlined'}
              onClick={() => setActiveTab('fees')}
              startIcon={<SettingsIcon />}
            >
              费率配置
            </Button>
          </Box>
        </Box>

        {/* 钱包管理 */}
        {activeTab === 'wallets' && (
          <Box>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
              <Typography variant="h6" component="h2">
                中间人钱包列表
              </Typography>
              <Button
                variant="contained"
                startIcon={<AddIcon />}
                onClick={() => setWalletDialogOpen(true)}
              >
                添加钱包
              </Button>
            </Box>

            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>钱包名称</TableCell>
                    <TableCell>钱包地址</TableCell>
                    <TableCell>余额 (USDT)</TableCell>
                    <TableCell>状态</TableCell>
                    <TableCell>创建时间</TableCell>
                    <TableCell>操作</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {wallets.map((wallet) => (
                    <TableRow key={wallet.id}>
                      <TableCell>{wallet.name}</TableCell>
                      <TableCell>
                        <Box sx={{ fontFamily: 'monospace', fontSize: '0.875rem' }}>
                          {wallet.wallet_address}
                        </Box>
                      </TableCell>
                      <TableCell>{wallet.balance?.toFixed(2) || '0.00'}</TableCell>
                      <TableCell>
                        <Chip
                          label={getStatusText(wallet.status)}
                          color={getStatusColor(wallet.status) as any}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>
                        {new Date(wallet.created_at).toLocaleString()}
                      </TableCell>
                      <TableCell>
                        <Box sx={{ display: 'flex', gap: 1 }}>
                          <FormControlLabel
                            control={
                              <Switch
                                checked={wallet.status === 'active'}
                                onChange={(e) => handleUpdateWalletStatus(wallet.id, e.target.checked ? 'active' : 'inactive')}
                                size="small"
                              />
                            }
                            label=""
                          />
                          <IconButton
                            size="small"
                            color="error"
                            onClick={() => handleDeleteWallet(wallet.id)}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Box>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Box>
        )}

        {/* 转账记录 */}
        {activeTab === 'transfers' && (
          <Box>
            <Typography variant="h6" component="h2" mb={3}>
              转账记录
            </Typography>

            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>转账ID</TableCell>
                    <TableCell>订单ID</TableCell>
                    <TableCell>来源地址</TableCell>
                    <TableCell>目标地址</TableCell>
                    <TableCell>转账金额</TableCell>
                    <TableCell>手续费</TableCell>
                    <TableCell>实际到账</TableCell>
                    <TableCell>状态</TableCell>
                    <TableCell>交易哈希</TableCell>
                    <TableCell>创建时间</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {transfers.map((transfer) => (
                    <TableRow key={transfer.id}>
                      <TableCell>{transfer.id}</TableCell>
                      <TableCell>{transfer.payment_order?.order_id || '-'}</TableCell>
                      <TableCell>
                        <Box sx={{ fontFamily: 'monospace', fontSize: '0.875rem' }}>
                          {transfer.from_address}
                        </Box>
                      </TableCell>
                      <TableCell>
                        <Box sx={{ fontFamily: 'monospace', fontSize: '0.875rem' }}>
                          {transfer.to_address}
                        </Box>
                      </TableCell>
                      <TableCell>{transfer.amount.toFixed(2)}</TableCell>
                      <TableCell>
                        <Box>
                          <div>固定: {transfer.fixed_fee_amount.toFixed(2)}</div>
                          <div>网络: {transfer.network_fee_amount.toFixed(2)}</div>
                          <div><strong>总计: {transfer.total_fee_amount.toFixed(2)}</strong></div>
                        </Box>
                      </TableCell>
                      <TableCell>{transfer.net_amount.toFixed(2)}</TableCell>
                      <TableCell>
                        <Chip
                          label={getStatusText(transfer.status)}
                          color={getStatusColor(transfer.status) as any}
                          size="small"
                        />
                        {transfer.error_message && (
                          <Typography variant="caption" display="block" color="error">
                            {transfer.error_message}
                          </Typography>
                        )}
                      </TableCell>
                      <TableCell>
                        {transfer.transaction_hash ? (
                          <Box sx={{ fontFamily: 'monospace', fontSize: '0.75rem' }}>
                            {transfer.transaction_hash}
                          </Box>
                        ) : (
                          '-'
                        )}
                      </TableCell>
                      <TableCell>
                        {new Date(transfer.created_at).toLocaleString()}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Box>
        )}

        {/* 费率配置 */}
        {activeTab === 'fees' && (
          <Box>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
              <Typography variant="h6" component="h2">
                手续费率配置
              </Typography>
              <Button
                variant="contained"
                startIcon={<SettingsIcon />}
                onClick={() => {
                                     if (feeConfig) {
                     setNewFeeConfig({
                       fee_rate: feeConfig.fee_rate,
                       min_fee: feeConfig.min_fee,
                       max_fee: feeConfig.max_fee,
                       network_fee_min: feeConfig.network_fee_min,
                       network_fee_max: feeConfig.network_fee_max
                     });
                   }
                  setFeeDialogOpen(true);
                }}
              >
                {feeConfig ? '编辑配置' : '创建配置'}
              </Button>
            </Box>

            {feeConfig && (
              <Card>
                <CardContent>
                  <Grid container spacing={3}>
                    <Grid item xs={12} md={2}>
                      <Typography variant="subtitle2" color="textSecondary">
                        固定费率
                      </Typography>
                      <Typography variant="h6">
                        {feeConfig.fee_rate}%
                      </Typography>
                    </Grid>
                    <Grid item xs={12} md={2}>
                      <Typography variant="subtitle2" color="textSecondary">
                        最小手续费
                      </Typography>
                      <Typography variant="h6">
                        {feeConfig.min_fee} USDT
                      </Typography>
                    </Grid>
                    <Grid item xs={12} md={2}>
                      <Typography variant="subtitle2" color="textSecondary">
                        最大手续费
                      </Typography>
                      <Typography variant="h6">
                        {feeConfig.max_fee} USDT
                      </Typography>
                    </Grid>
                    <Grid item xs={12} md={2}>
                      <Typography variant="subtitle2" color="textSecondary">
                        网络手续费范围
                      </Typography>
                      <Typography variant="h6">
                        {feeConfig.network_fee_min}-{feeConfig.network_fee_max} USDT
                      </Typography>
                    </Grid>
                                         <Grid item xs={12} md={2}>
                       <Typography variant="subtitle2" color="textSecondary">
                         网络转账次数
                       </Typography>
                       <Typography variant="h6">
                         2 次
                       </Typography>
                     </Grid>
                  </Grid>
                  <Box mt={2}>
                    <Typography variant="body2" color="textSecondary">
                      最后更新: {new Date(feeConfig.updated_at).toLocaleString()}
                    </Typography>
                  </Box>
                </CardContent>
              </Card>
            )}
          </Box>
        )}

        {/* 添加钱包对话框 */}
        <Dialog open={walletDialogOpen} onClose={() => setWalletDialogOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>添加中间人钱包</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="钱包名称"
              fullWidth
              variant="outlined"
              value={newWallet.name}
              onChange={(e) => setNewWallet({ ...newWallet, name: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="钱包地址"
              fullWidth
              variant="outlined"
              value={newWallet.wallet_address}
              onChange={(e) => setNewWallet({ ...newWallet, wallet_address: e.target.value })}
              placeholder="TRC20钱包地址"
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="私钥"
              fullWidth
              variant="outlined"
              type="password"
              value={newWallet.private_key}
              onChange={(e) => setNewWallet({ ...newWallet, private_key: e.target.value })}
              placeholder="钱包私钥（用于转账）"
              helperText="⚠️ 私钥用于USDT转账，请妥善保管，不要泄露给他人"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setWalletDialogOpen(false)}>取消</Button>
            <Button onClick={handleAddWallet} variant="contained">添加</Button>
          </DialogActions>
        </Dialog>

        {/* 费率配置对话框 */}
        <Dialog open={feeDialogOpen} onClose={() => setFeeDialogOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>
            {feeConfig ? '编辑费率配置' : '创建费率配置'}
          </DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="手续费率 (%)"
              type="number"
              fullWidth
              variant="outlined"
              value={newFeeConfig.fee_rate}
              onChange={(e) => setNewFeeConfig({ ...newFeeConfig, fee_rate: parseFloat(e.target.value) || 0 })}
              sx={{ mb: 2 }}
              inputProps={{ step: 0.1, min: 0, max: 100 }}
            />
            <TextField
              margin="dense"
              label="最小手续费 (USDT)"
              type="number"
              fullWidth
              variant="outlined"
              value={newFeeConfig.min_fee}
              onChange={(e) => setNewFeeConfig({ ...newFeeConfig, min_fee: parseFloat(e.target.value) || 0 })}
              sx={{ mb: 2 }}
              inputProps={{ step: 0.1, min: 0 }}
            />
            <TextField
              margin="dense"
              label="最大手续费 (USDT)"
              type="number"
              fullWidth
              variant="outlined"
              value={newFeeConfig.max_fee}
              onChange={(e) => setNewFeeConfig({ ...newFeeConfig, max_fee: parseFloat(e.target.value) || 0 })}
              inputProps={{ step: 0.1, min: 0 }}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="网络手续费最小值 (USDT)"
              type="number"
              fullWidth
              variant="outlined"
              value={newFeeConfig.network_fee_min}
              onChange={(e) => setNewFeeConfig({ ...newFeeConfig, network_fee_min: parseFloat(e.target.value) || 0 })}
              inputProps={{ step: 0.001, min: 0 }}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="网络手续费最大值 (USDT)"
              type="number"
              fullWidth
              variant="outlined"
              value={newFeeConfig.network_fee_max}
              onChange={(e) => setNewFeeConfig({ ...newFeeConfig, network_fee_max: parseFloat(e.target.value) || 0 })}
              inputProps={{ step: 0.001, min: 0 }}
              sx={{ mb: 2 }}
            />
            
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setFeeDialogOpen(false)}>取消</Button>
            <Button onClick={handleUpdateFeeConfig} variant="contained">
              {feeConfig ? '更新' : '创建'}
            </Button>
          </DialogActions>
        </Dialog>
      </Container>
    </>
  );
};

export default AdminMiddlemanPage;
