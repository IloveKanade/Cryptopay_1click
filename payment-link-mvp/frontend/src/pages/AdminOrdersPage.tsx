import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Typography,
  Box,
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Chip,
  Alert,
  CircularProgress,
  InputAdornment,
  Tooltip,
  Menu,
  MenuItem,
  ListItemIcon,
  ListItemText,
  Card,
  CardContent,
  Grid,
} from '@mui/material';
import AdminNavigation from '../components/AdminNavigation';
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  MoreVert as MoreVertIcon,
  Receipt as ReceiptIcon,
  Refresh as RefreshIcon,
  Visibility as ViewIcon,
  Payment as PaymentIcon,
  Schedule as PendingIcon,
  CheckCircle as CompletedIcon,
  Cancel as FailedIcon,
} from '@mui/icons-material';
import { api } from '../services/api';

interface Order {
  id: number;
  trade_id: string;
  order_id: string;
  amount: number;
  actual_amount: number;
  status: string;
  created_at: string;
  updated_at: string;
  block_transaction_id?: string;
  payment_link?: {
    id: number;
    title: string;
    description: string;
    amount: number;
    user: {
      id: number;
      name: string;
      email: string;
    };
  };
}

interface OrderListResponse {
  orders: Order[];
  total: number;
}

const AdminOrdersPage: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchQuery, setSearchQuery] = useState('');
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
  const [editForm, setEditForm] = useState({
    status: '',
  });
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedOrderId, setSelectedOrderId] = useState<number | null>(null);
  const [viewDialogOpen, setViewDialogOpen] = useState(false);
  const [viewOrder, setViewOrder] = useState<Order | null>(null);

  useEffect(() => {
    fetchOrders();
  }, [page, rowsPerPage, searchQuery]);

  const fetchOrders = async () => {
    try {
      setLoading(true);
      const params = {
        page: page + 1,
        page_size: rowsPerPage,
      };

      let response;
      if (searchQuery) {
        response = await api.get('/api/admin/orders/search', { params: { ...params, q: searchQuery } });
      } else {
        response = await api.get('/api/admin/orders', { params });
      }

      const data: OrderListResponse = response.data.data;
      setOrders(data.orders);
      setTotal(data.total);
      setError(null);
    } catch (err: any) {
      setError(err.response?.data?.error || '获取订单列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleEditOrder = (order: Order) => {
    setSelectedOrder(order);
    setEditForm({
      status: order.status,
    });
    setEditDialogOpen(true);
  };

  const handleUpdateOrder = async () => {
    if (!selectedOrder) return;

    try {
              await api.put(`/api/admin/orders/${selectedOrder.id}/status`, editForm);
      setEditDialogOpen(false);
      fetchOrders();
    } catch (err: any) {
      setError(err.response?.data?.error || '更新订单失败');
    }
  };

  const handleDeleteOrder = async (orderId: number) => {
    if (!window.confirm('确定要删除这个订单吗？此操作不可撤销。')) {
      return;
    }

    try {
              await api.delete(`/api/admin/orders/${orderId}`);
      fetchOrders();
    } catch (err: any) {
      setError(err.response?.data?.error || '删除订单失败');
    }
  };

  const handleViewOrder = async (orderId: number) => {
    try {
              const response = await api.get(`/api/admin/orders/${orderId}`);
      setViewOrder(response.data.data);
      setViewDialogOpen(true);
    } catch (err: any) {
      setError(err.response?.data?.error || '获取订单详情失败');
    }
  };

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, orderId: number) => {
    setAnchorEl(event.currentTarget);
    setSelectedOrderId(orderId);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedOrderId(null);
  };

  const getStatusChip = (status: string) => {
    switch (status) {
      case 'pending':
        return <Chip label="待支付" color="warning" size="small" icon={<PendingIcon />} />;
      case 'paid':
        return <Chip label="已支付" color="success" size="small" icon={<CompletedIcon />} />;
      case 'expired':
        return <Chip label="已过期" color="error" size="small" icon={<FailedIcon />} />;
      default:
        return <Chip label={status} color="default" size="small" />;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('zh-CN');
  };

  const formatAmount = (amount: number) => {
    return `$${amount.toFixed(2)}`;
  };

  return (
    <>
      <AdminNavigation />
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        {/* 页面标题 */}
        <Box display="flex" alignItems="center" justifyContent="space-between" mb={4}>
        <Box display="flex" alignItems="center">
          <ReceiptIcon sx={{ mr: 2, fontSize: 32, color: 'primary.main' }} />
          <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold' }}>
            订单管理
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<RefreshIcon />}
          onClick={fetchOrders}
          disabled={loading}
        >
          刷新
        </Button>
      </Box>

      {/* 搜索栏 */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="搜索订单ID或交易ID..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
          }}
        />
      </Paper>

      {/* 错误提示 */}
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {/* 订单表格 */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>交易ID</TableCell>
                <TableCell>订单ID</TableCell>
                <TableCell>金额</TableCell>
                <TableCell>状态</TableCell>
                <TableCell>创建时间</TableCell>
                <TableCell>用户</TableCell>
                <TableCell align="right">操作</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={8} align="center">
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : orders.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} align="center">
                    <Typography variant="body2" color="textSecondary">
                      暂无订单数据
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                orders.map((order) => (
                  <TableRow key={order.id} hover>
                    <TableCell>{order.id}</TableCell>
                    <TableCell>
                      <Typography variant="body2" fontFamily="monospace">
                        {order.trade_id}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" fontFamily="monospace">
                        {order.order_id}
                      </Typography>
                    </TableCell>
                    <TableCell>{formatAmount(order.amount)}</TableCell>
                    <TableCell>{getStatusChip(order.status)}</TableCell>
                    <TableCell>{formatDate(order.created_at)}</TableCell>
                    <TableCell>
                      {order.payment_link?.user ? (
                        <Box>
                          <Typography variant="body2">{order.payment_link.user.name}</Typography>
                          <Typography variant="caption" color="textSecondary">
                            {order.payment_link.user.email}
                          </Typography>
                        </Box>
                      ) : (
                        <Typography variant="body2" color="textSecondary">
                          未知用户
                        </Typography>
                      )}
                    </TableCell>
                    <TableCell align="right">
                      <Tooltip title="更多操作">
                        <IconButton
                          onClick={(e) => handleMenuOpen(e, order.id)}
                          size="small"
                        >
                          <MoreVertIcon />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>

        {/* 分页 */}
        <TablePagination
          component="div"
          count={total}
          page={page}
          onPageChange={(_, newPage) => setPage(newPage)}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={(e) => {
            setRowsPerPage(parseInt(e.target.value, 10));
            setPage(0);
          }}
          labelRowsPerPage="每页行数:"
          labelDisplayedRows={({ from, to, count }) =>
            `${from}-${to} / ${count !== -1 ? count : `超过 ${to}`}`
          }
        />
      </Paper>

      {/* 操作菜单 */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem
          onClick={() => {
            if (selectedOrderId) {
              handleViewOrder(selectedOrderId);
            }
            handleMenuClose();
          }}
        >
          <ListItemIcon>
            <ViewIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>查看详情</ListItemText>
        </MenuItem>
        <MenuItem
          onClick={() => {
            const order = orders.find(o => o.id === selectedOrderId);
            if (order) {
              handleEditOrder(order);
            }
            handleMenuClose();
          }}
        >
          <ListItemIcon>
            <EditIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>编辑状态</ListItemText>
        </MenuItem>
        <MenuItem
          onClick={() => {
            if (selectedOrderId) {
              handleDeleteOrder(selectedOrderId);
            }
            handleMenuClose();
          }}
          sx={{ color: 'error.main' }}
        >
          <ListItemIcon>
            <DeleteIcon fontSize="small" color="error" />
          </ListItemIcon>
          <ListItemText>删除订单</ListItemText>
        </MenuItem>
      </Menu>

      {/* 编辑订单对话框 */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>编辑订单状态</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            select
            label="状态"
            value={editForm.status}
            onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}
            margin="normal"
          >
            <MenuItem value="pending">待支付</MenuItem>
            <MenuItem value="paid">已支付</MenuItem>
            <MenuItem value="expired">已过期</MenuItem>
          </TextField>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>取消</Button>
          <Button onClick={handleUpdateOrder} variant="contained">
            保存
          </Button>
        </DialogActions>
      </Dialog>

      {/* 查看订单详情对话框 */}
      <Dialog open={viewDialogOpen} onClose={() => setViewDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>订单详情</DialogTitle>
        <DialogContent>
          {viewOrder && (
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      订单信息
                    </Typography>
                    <Box>
                      <Typography variant="body2" color="textSecondary">
                        订单ID: {viewOrder.order_id}
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        交易ID: {viewOrder.trade_id}
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        金额: {formatAmount(viewOrder.amount)}
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        实际金额: {formatAmount(viewOrder.actual_amount)}
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        状态: {viewOrder.status}
                      </Typography>
                      {viewOrder.block_transaction_id && (
                        <Typography variant="body2" color="textSecondary">
                          交易哈希: {viewOrder.block_transaction_id}
                        </Typography>
                      )}
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} md={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      用户信息
                    </Typography>
                    {viewOrder.payment_link?.user ? (
                      <Box>
                        <Typography variant="body2" color="textSecondary">
                          用户ID: {viewOrder.payment_link.user.id}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          姓名: {viewOrder.payment_link.user.name}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                          邮箱: {viewOrder.payment_link.user.email}
                        </Typography>
                      </Box>
                    ) : (
                      <Typography variant="body2" color="textSecondary">
                        未知用户
                      </Typography>
                    )}
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      时间信息
                    </Typography>
                    <Box>
                      <Typography variant="body2" color="textSecondary">
                        创建时间: {formatDate(viewOrder.created_at)}
                      </Typography>
                      <Typography variant="body2" color="textSecondary">
                        更新时间: {formatDate(viewOrder.updated_at)}
                      </Typography>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setViewDialogOpen(false)}>关闭</Button>
        </DialogActions>
      </Dialog>
      </Container>
    </>
  );
};

export default AdminOrdersPage;
