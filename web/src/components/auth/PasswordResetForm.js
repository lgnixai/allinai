import React, { useEffect, useState } from 'react';
import { API, getLogo, showError, showInfo, showSuccess, getSystemName } from '../../helpers';
import Turnstile from 'react-turnstile';
import { Button, Card, Form, Typography } from '@douyinfe/semi-ui';
import { IconPhone, IconKey, IconLock } from '@douyinfe/semi-icons';
import { Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const { Text, Title } = Typography;

const PasswordResetForm = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [inputs, setInputs] = useState({
    phone: '',
    verification_code: '',
  });
  const { phone, verification_code } = inputs;

  const [loading, setLoading] = useState(false);
  const [turnstileEnabled, setTurnstileEnabled] = useState(false);
  const [turnstileSiteKey, setTurnstileSiteKey] = useState('');
  const [turnstileToken, setTurnstileToken] = useState('');
  const [disableButton, setDisableButton] = useState(false);
  const [countdown, setCountdown] = useState(30);
  const [verificationCodeLoading, setVerificationCodeLoading] = useState(false);
  const [showPasswordForm, setShowPasswordForm] = useState(false);
  const [passwordForm, setPasswordForm] = useState({
    password: '',
    password2: '',
  });
  const { password, password2 } = passwordForm;

  const logo = getLogo();
  const systemName = getSystemName();

  useEffect(() => {
    let status = localStorage.getItem('status');
    if (status) {
      status = JSON.parse(status);
      if (status.turnstile_check) {
        setTurnstileEnabled(true);
        setTurnstileSiteKey(status.turnstile_site_key);
      }
    }
  }, []);

  useEffect(() => {
    let countdownInterval = null;
    if (disableButton && countdown > 0) {
      countdownInterval = setInterval(() => {
        setCountdown(countdown - 1);
      }, 1000);
    } else if (countdown === 0) {
      setDisableButton(false);
      setCountdown(30);
    }
    return () => clearInterval(countdownInterval);
  }, [disableButton, countdown]);

  function handleChange(name, value) {
    setInputs((inputs) => ({ ...inputs, [name]: value }));
  }

  function handlePasswordChange(name, value) {
    setPasswordForm((form) => ({ ...form, [name]: value }));
  }

  const sendVerificationCode = async () => {
    if (phone === '') {
      showInfo(t('请先输入手机号！'));
      return;
    }
    if (phone.length !== 11) {
      showInfo(t('请输入正确的11位手机号！'));
      return;
    }
    if (turnstileEnabled && turnstileToken === '') {
      showInfo(t('请稍后几秒重试，Turnstile 正在检查用户环境！'));
      return;
    }
    setVerificationCodeLoading(true);
    try {
      const res = await API.get(
        `/api/reset_password?phone=${phone}&turnstile=${turnstileToken}`,
      );
      const { success, message } = res.data;
      if (success) {
        showSuccess(t('重置验证码发送成功，请检查手机！'));
        setDisableButton(true);
      } else {
        showError(message);
      }
    } catch (error) {
      showError('发送验证码失败，请重试');
    } finally {
      setVerificationCodeLoading(false);
    }
  };

  async function handleVerifyCode() {
    if (!phone) {
      showError(t('请输入手机号'));
      return;
    }
    if (phone.length !== 11) {
      showError(t('请输入正确的11位手机号'));
      return;
    }
    if (!verification_code) {
      showError(t('请输入验证码'));
      return;
    }
    if (turnstileEnabled && turnstileToken === '') {
      showInfo(t('请稍后几秒重试，Turnstile 正在检查用户环境！'));
      return;
    }
    setLoading(true);
    try {
      const res = await API.post(`/api/user/verify_reset_code`, {
        phone,
        token: verification_code,
      });
      const { success, message } = res.data;
      if (success) {
        showSuccess(t('验证码验证成功！'));
        setShowPasswordForm(true);
      } else {
        showError(message);
      }
    } catch (error) {
      showError('验证码验证失败，请重试');
    } finally {
      setLoading(false);
    }
  }

  async function handleResetPassword() {
    if (!password) {
      showError(t('请输入密码'));
      return;
    }
    if (password.length < 8) {
      showError(t('密码长度不能少于8位'));
      return;
    }
    if (password.length > 20) {
      showError(t('密码长度不能超过20位'));
      return;
    }
    if (password !== password2) {
      showError(t('两次输入的密码不一致'));
      return;
    }
    if (turnstileEnabled && turnstileToken === '') {
      showInfo(t('请稍后几秒重试，Turnstile 正在检查用户环境！'));
      return;
    }
    setLoading(true);
    try {
      const res = await API.post(`/api/user/reset_password`, {
        phone,
        token: verification_code,
        password,
      });
      const { success, message } = res.data;
      if (success) {
        showSuccess(t('密码重置成功！'));
        navigate('/login');
      } else {
        showError(message);
      }
    } catch (error) {
      showError('密码重置失败，请重试');
    } finally {
      setLoading(false);
    }
  }

  const renderVerificationForm = () => {
    return (
      <Form className="space-y-3">
        <Form.Input
          field="phone"
          label={t('手机号')}
          placeholder={t('请输入您的手机号')}
          name="phone"
          size="large"
          value={phone}
          onChange={(value) => handleChange('phone', value)}
          prefix={<IconPhone />}
          maxLength={11}
        />

        <Form.Input
          field="verification_code"
          label={t('手机验证码')}
          placeholder={t('输入验证码')}
          name="verification_code"
          size="large"
          value={verification_code}
          onChange={(value) => handleChange('verification_code', value)}
          prefix={<IconKey />}
          suffix={
            <Button
              onClick={sendVerificationCode}
              loading={verificationCodeLoading}
              disabled={disableButton}
              size="small"
            >
              {disableButton ? `${t('重试')} (${countdown})` : t('获取验证码')}
            </Button>
          }
        />

        <div className="space-y-2 pt-2">
          <Button
            theme="solid"
            className="w-full !rounded-full"
            type="primary"
            htmlType="submit"
            size="large"
            onClick={handleVerifyCode}
            loading={loading}
          >
            {t('验证验证码')}
          </Button>
        </div>
      </Form>
    );
  };

  const renderPasswordForm = () => {
    return (
      <Form className="space-y-3">
        <Form.Input
          field="password"
          label={t('新密码')}
          placeholder={t('输入密码，最短 8 位，最长 20 位')}
          name="password"
          mode="password"
          size="large"
          value={password}
          onChange={(value) => handlePasswordChange('password', value)}
          prefix={<IconLock />}
        />

        <Form.Input
          field="password2"
          label={t('确认密码')}
          placeholder={t('确认密码')}
          name="password2"
          mode="password"
          size="large"
          value={password2}
          onChange={(value) => handlePasswordChange('password2', value)}
          prefix={<IconLock />}
        />

        <div className="space-y-2 pt-2">
          <Button
            theme="solid"
            className="w-full !rounded-full"
            type="primary"
            htmlType="submit"
            size="large"
            onClick={handleResetPassword}
            loading={loading}
          >
            {t('重置密码')}
          </Button>
        </div>
      </Form>
    );
  };

  return (
    <div className="relative overflow-hidden bg-gray-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      {/* 背景模糊晕染球 */}
      <div className="blur-ball blur-ball-indigo" style={{ top: '-80px', right: '-80px', transform: 'none' }} />
      <div className="blur-ball blur-ball-teal" style={{ top: '50%', left: '-120px' }} />
      <div className="w-full max-w-sm mt-[64px]">
        <div className="flex flex-col items-center">
          <div className="w-full max-w-md">
            <div className="flex items-center justify-center mb-6 gap-2">
              <img src={logo} alt="Logo" className="h-10 rounded-full" />
              <Title heading={3} className='!text-gray-800'>{systemName}</Title>
            </div>

            <Card className="shadow-xl border-0 !rounded-2xl overflow-hidden">
              <div className="flex justify-center pt-6 pb-2">
                <Title heading={3} className="text-gray-800 dark:text-gray-200">
                  {showPasswordForm ? t('设置新密码') : t('密码重置')}
                </Title>
              </div>
              <div className="px-2 py-8">
                {showPasswordForm ? renderPasswordForm() : renderVerificationForm()}

                <div className="mt-6 text-center text-sm">
                  <Text>{t('想起来了？')} <Link to="/login" className="text-blue-600 hover:text-blue-800 font-medium">{t('登录')}</Link></Text>
                </div>
              </div>
            </Card>

            {turnstileEnabled && (
              <div className="flex justify-center mt-6">
                <Turnstile
                  sitekey={turnstileSiteKey}
                  onVerify={(token) => {
                    setTurnstileToken(token);
                  }}
                />
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default PasswordResetForm;
