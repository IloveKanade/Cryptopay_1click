import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Box,
  Button,
  Card,
  CardContent,
  Grid,
  TextField,
  Alert,
  CircularProgress,
  Divider,
  Paper,
  Switch,
  FormControlLabel,
} from '@mui/material';
import {
  Settings as SettingsIcon,
  Key as KeyIcon,
  Security as SecurityIcon,
  Save as SaveIcon,
} from '@mui/icons-material';
import { api } from '../services/api';
import AdminNavigation from '../components/AdminNavigation';

interface SystemSettings {
  epusdt_base_url: string;
  callback_success_url: string;
  callback_fail_url: string;
  notify_secret: string;
}

const AdminSettingsPage: React.FC = () => {
  const [settings, setSettings] = useState<SystemSettings>({
    epusdt_base_url: '',
    callback_success_url: '',
    callback_fail_url: '',
    notify_secret: '',
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // 获取Epusdt配置
      const epusdtResponse = await api.get('/api/admin/epusdt/config');
      const callbackResponse = await api.get('/api/admin/callback/config');
      
      setSettings({
        epusdt_base_url: epusdtResponse.data.data?.base_url || '',
        callback_success_url: callbackResponse.data.data?.transfer_success_url || '',
        callback_fail_url: callbackResponse.data.data?.transfer_fail_url || '',
        notify_secret: callbackResponse.data.data?.notify_secret || '',
      });
    } catch (err: any) {
      setError(err.response?.data?.error || '获取设置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSaveSettings = async () => {
    try {
      setSaving(true);
      setError(null);
      setSuccess(null);

      // 保存Epusdt配置
      await api.put('/api/admin/epusdt/config', {
        base_url: settings.epusdt_base_url,
      });

      // 保存回调配置
      await api.put('/api/admin/callback/config', {
        transfer_success_url: settings.callback_success_url,
        transfer_fail_url: settings.callback_fail_url,
        notify_secret: settings.notify_secret,
      });

      setSuccess('设置保存成功！');
    } catch (err: any) {
      setError(err.response?.data?.error || '保存设置失败');
    } finally {
      setSaving(false);
    }
  };

  const handleTestEpusdtConnection = async () => {
    try {
      setError(null);
      await api.post('/api/admin/epusdt/test-connection');
      setSuccess('Epusdt连接测试成功！');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Epusdt连接测试失败');
    }
  };

  const handleTestCallback = async () => {
    try {
      setError(null);
      await api.post('/api/admin/callback/test');
      setSuccess('回调测试成功！');
    } catch (err: any) {
      setError(err.response?.data?.error || '回调测试失败');
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
          <SettingsIcon sx={{ mr: 2, fontSize: 32, color: 'primary.main' }} />
          <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold' }}>
            系统设置
          </Typography>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        {success && (
          <Alert severity="success" sx={{ mb: 3 }}>
            {success}
          </Alert>
        )}

        <Grid container spacing={3}>
          {/* Epusdt配置 */}
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" mb={3}>
                  <KeyIcon sx={{ mr: 1, color: 'primary.main' }} />
                  <Typography variant="h6" component="h2">
                    Epusdt配置
                  </Typography>
                </Box>
                <Divider sx={{ mb: 3 }} />
                
                <TextField
                  fullWidth
                  label="Epusdt基础URL"
                  variant="outlined"
                  value={settings.epusdt_base_url}
                  onChange={(e) => setSettings({ ...settings, epusdt_base_url: e.target.value })}
                  placeholder="http://localhost:8000"
                  sx={{ mb: 2 }}
                  helperText="Epusdt服务的访问地址"
                />
                

                
                <Box display="flex" gap={2}>
                  <Button
                    variant="outlined"
                    onClick={handleTestEpusdtConnection}
                    disabled={saving}
                  >
                    测试连接
                  </Button>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* 回调配置 */}
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" mb={3}>
                  <SecurityIcon sx={{ mr: 1, color: 'primary.main' }} />
                  <Typography variant="h6" component="h2">
                    回调配置
                  </Typography>
                </Box>
                <Divider sx={{ mb: 3 }} />
                
                <TextField
                  fullWidth
                  label="转账成功回调地址"
                  variant="outlined"
                  value={settings.callback_success_url}
                  onChange={(e) => setSettings({ ...settings, callback_success_url: e.target.value })}
                  placeholder="https://your-domain.com/callback/success"
                  sx={{ mb: 2 }}
                  helperText="转账成功时的回调通知地址"
                />
                
                <TextField
                  fullWidth
                  label="转账失败回调地址"
                  variant="outlined"
                  value={settings.callback_fail_url}
                  onChange={(e) => setSettings({ ...settings, callback_fail_url: e.target.value })}
                  placeholder="https://your-domain.com/callback/fail"
                  sx={{ mb: 2 }}
                  helperText="转账失败时的回调通知地址"
                />
                
                <TextField
                  fullWidth
                  label="回调签名密钥"
                  variant="outlined"
                  type="password"
                  value={settings.notify_secret}
                  onChange={(e) => setSettings({ ...settings, notify_secret: e.target.value })}
                  placeholder="callback_secret_key"
                  sx={{ mb: 2 }}
                  helperText="用于验证回调请求的真实性，防止伪造请求"
                />
                
                <Box display="flex" gap={2}>
                  <Button
                    variant="outlined"
                    onClick={handleTestCallback}
                    disabled={saving}
                  >
                    测试回调
                  </Button>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* 保存按钮 */}
        <Box display="flex" justifyContent="center" mt={4}>
          <Button
            variant="contained"
            size="large"
            startIcon={<SaveIcon />}
            onClick={handleSaveSettings}
            disabled={saving}
            sx={{ minWidth: 200 }}
          >
            {saving ? '保存中...' : '保存设置'}
          </Button>
        </Box>

        {/* 说明信息 */}
        <Paper sx={{ p: 3, mt: 4 }}>
          <Typography variant="h6" gutterBottom>
            配置说明
          </Typography>
          <Divider sx={{ mb: 2 }} />
          <Typography variant="body2" color="textSecondary" paragraph>
            <strong>Epusdt配置：</strong>
          </Typography>
                     <Typography variant="body2" color="textSecondary" paragraph>
             • <strong>基础URL：</strong>Epusdt服务的访问地址，例如：http://localhost:8000
           </Typography>
          <Typography variant="body2" color="textSecondary" paragraph>
            <strong>回调配置：</strong>
          </Typography>
          <Typography variant="body2" color="textSecondary" paragraph>
            • <strong>成功回调地址：</strong>转账成功时系统会向此地址发送通知
          </Typography>
          <Typography variant="body2" color="textSecondary" paragraph>
            • <strong>失败回调地址：</strong>转账失败时系统会向此地址发送通知
          </Typography>
          <Typography variant="body2" color="textSecondary" paragraph>
            • <strong>回调签名密钥：</strong>用于验证回调请求的真实性，防止伪造请求
          </Typography>
        </Paper>
      </Container>
    </>
  );
};

export default AdminSettingsPage;
