import React, { useContext, useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  API,
  getLogo,
  showError,
  showInfo,
  showSuccess,
  updateAPI,
  getSystemName,
  setUserData
} from '../../helpers/index.js';
import Turnstile from 'react-turnstile';
import {
  Button,
  Card,
  Divider,
  Form,
  Icon,
  Modal,
} from '@douyinfe/semi-ui';
import Title from '@douyinfe/semi-ui/lib/es/typography/title';
import Text from '@douyinfe/semi-ui/lib/es/typography/text';
import { IconGithubLogo, IconMail, IconUser, IconLock, IconKey, IconPhone } from '@douyinfe/semi-icons';
import {
  onGitHubOAuthClicked,
  onLinuxDOOAuthClicked,
  onOIDCClicked,
} from '../../helpers/index.js';
import OIDCIcon from '../common/logo/OIDCIcon.js';
import LinuxDoIcon from '../common/logo/LinuxDoIcon.js';
import WeChatIcon from '../common/logo/WeChatIcon.js';
import TelegramLoginButton from 'react-telegram-login/src';
import { UserContext } from '../../context/User/index.js';
import { useTranslation } from 'react-i18next';

const RegisterForm = () => {
  let navigate = useNavigate();
  const { t } = useTranslation();
  const [inputs, setInputs] = useState({
    phone: '',
    display_name: '',
    school: '',
    college: '',
    phone_verification_code: '',
    wechat_verification_code: '',
  });
  const { phone, display_name, school, college } = inputs;
  const [userState, userDispatch] = useContext(UserContext);
  const [turnstileEnabled, setTurnstileEnabled] = useState(false);
  const [turnstileSiteKey, setTurnstileSiteKey] = useState('');
  const [turnstileToken, setTurnstileToken] = useState('');
  const [showWeChatLoginModal, setShowWeChatLoginModal] = useState(false);
  const [showPhoneRegister, setShowPhoneRegister] = useState(false);
  const [wechatLoading, setWechatLoading] = useState(false);
  const [githubLoading, setGithubLoading] = useState(false);
  const [oidcLoading, setOidcLoading] = useState(false);
  const [linuxdoLoading, setLinuxdoLoading] = useState(false);
  const [phoneRegisterLoading, setPhoneRegisterLoading] = useState(false);
  const [registerLoading, setRegisterLoading] = useState(false);
  const [verificationCodeLoading, setVerificationCodeLoading] = useState(false);
  const [otherRegisterOptionsLoading, setOtherRegisterOptionsLoading] = useState(false);
  const [wechatCodeSubmitLoading, setWechatCodeSubmitLoading] = useState(false);

  const logo = getLogo();
  const systemName = getSystemName();

  let affCode = new URLSearchParams(window.location.search).get('aff');
  if (affCode) {
    localStorage.setItem('aff', affCode);
  }

  const [status] = useState(() => {
    const savedStatus = localStorage.getItem('status');
    return savedStatus ? JSON.parse(savedStatus) : {};
  });

  useEffect(() => {
    if (status.turnstile_check) {
      setTurnstileEnabled(true);
      setTurnstileSiteKey(status.turnstile_site_key);
    }
  }, [status]);

  const onWeChatLoginClicked = () => {
    setWechatLoading(true);
    setShowWeChatLoginModal(true);
    setWechatLoading(false);
  };

  const onSubmitWeChatVerificationCode = async () => {
    if (turnstileEnabled && turnstileToken === '') {
      showInfo('请稍后几秒重试，Turnstile 正在检查用户环境！');
      return;
    }
    setWechatCodeSubmitLoading(true);
    try {
      const res = await API.get(
        `/api/oauth/wechat?code=${inputs.wechat_verification_code}`,
      );
      const { success, message, data } = res.data;
      if (success) {
        userDispatch({ type: 'login', payload: data });
        localStorage.setItem('user', JSON.stringify(data));
        setUserData(data);
        updateAPI();
        navigate('/');
        showSuccess('登录成功！');
        setShowWeChatLoginModal(false);
      } else {
        showError(message);
      }
    } catch (error) {
      showError('登录失败，请重试');
    } finally {
      setWechatCodeSubmitLoading(false);
    }
  };

  function handleChange(name, value) {
    setInputs((inputs) => ({ ...inputs, [name]: value }));
  }

  async function handleSubmit(e) {
    if (!phone || phone.length !== 11) {
      showInfo('请输入正确的11位手机号！');
      return;
    }
    if (!inputs.phone_verification_code) {
      showInfo('请输入手机验证码！');
      return;
    }
    if (phone) {
      if (turnstileEnabled && turnstileToken === '') {
        showInfo('请稍后几秒重试，Turnstile 正在检查用户环境！');
        return;
      }
      setRegisterLoading(true);
      try {
        if (!affCode) {
          affCode = localStorage.getItem('aff');
        }
        inputs.aff_code = affCode;
        const res = await API.post(
          `/api/user/register?turnstile=${turnstileToken}`,
          inputs,
        );
        const { success, message } = res.data;
        if (success) {
          navigate('/login');
          showSuccess('注册成功！');
        } else {
          showError(message);
        }
      } catch (error) {
        showError('注册失败，请重试');
      } finally {
        setRegisterLoading(false);
      }
    }
  }

  const sendVerificationCode = async () => {
    if (inputs.phone === '') {
      showInfo('请先输入手机号！');
      return;
    }
    if (inputs.phone.length !== 11) {
      showInfo('请输入正确的11位手机号！');
      return;
    }
    if (turnstileEnabled && turnstileToken === '') {
      showInfo('请稍后几秒重试，Turnstile 正在检查用户环境！');
      return;
    }
    setVerificationCodeLoading(true);
    try {
      const res = await API.get(
        `/api/phone_verification?phone=${inputs.phone}&purpose=register&turnstile=${turnstileToken}`,
      );
      const { success, message } = res.data;
      if (success) {
        showSuccess('验证码发送成功，请检查你的手机！');
      } else {
        showError(message);
      }
    } catch (error) {
      showError('发送验证码失败，请重试');
    } finally {
      setVerificationCodeLoading(false);
    }
  };

  const handleGitHubClick = () => {
    setGithubLoading(true);
    try {
      onGitHubOAuthClicked(status.github_client_id);
    } finally {
      setTimeout(() => setGithubLoading(false), 3000);
    }
  };

  const handleOIDCClick = () => {
    setOidcLoading(true);
    try {
      onOIDCClicked(
        status.oidc_authorization_endpoint,
        status.oidc_client_id
      );
    } finally {
      setTimeout(() => setOidcLoading(false), 3000);
    }
  };

  const handleLinuxDOClick = () => {
    setLinuxdoLoading(true);
    try {
      onLinuxDOOAuthClicked(status.linuxdo_client_id);
    } finally {
      setTimeout(() => setLinuxdoLoading(false), 3000);
    }
  };

  const handlePhoneRegisterClick = () => {
    setPhoneRegisterLoading(true);
    setShowPhoneRegister(true);
    setPhoneRegisterLoading(false);
  };

  const handleOtherRegisterOptionsClick = () => {
    setOtherRegisterOptionsLoading(true);
    setShowPhoneRegister(false);
    setOtherRegisterOptionsLoading(false);
  };

  const onTelegramLoginClicked = async (response) => {
    const fields = [
      'id',
      'first_name',
      'last_name',
      'username',
      'photo_url',
      'auth_date',
      'hash',
      'lang',
    ];
    const params = {};
    fields.forEach((field) => {
      if (response[field]) {
        params[field] = response[field];
      }
    });
    try {
      const res = await API.get(`/api/oauth/telegram/login`, { params });
      const { success, message, data } = res.data;
      if (success) {
        userDispatch({ type: 'login', payload: data });
        localStorage.setItem('user', JSON.stringify(data));
        showSuccess('登录成功！');
        setUserData(data);
        updateAPI();
        navigate('/');
      } else {
        showError(message);
      }
    } catch (error) {
      showError('登录失败，请重试');
    }
  };

  const renderOAuthOptions = () => {
    return (
      <div className="flex flex-col items-center">
        <div className="w-full max-w-md">
          <div className="flex items-center justify-center mb-6 gap-2">
            <img src={logo} alt="Logo" className="h-10 rounded-full" />
            <Title heading={3} className='!text-gray-800'>{systemName}</Title>
          </div>

          <Card className="shadow-xl border-0 !rounded-2xl overflow-hidden">
            <div className="flex justify-center pt-6 pb-2">
              <Title heading={3} className="text-gray-800 dark:text-gray-200">{t('注 册')}</Title>
            </div>
            <div className="px-2 py-8">
              <div className="space-y-3">
                {status.wechat_login && (
                  <Button
                    theme='outline'
                    className="w-full h-12 flex items-center justify-center !rounded-full border border-gray-200 hover:bg-gray-50 transition-colors"
                    type="tertiary"
                    icon={<Icon svg={<WeChatIcon />} style={{ color: '#07C160' }} />}
                    size="large"
                    onClick={onWeChatLoginClicked}
                    loading={wechatLoading}
                  >
                    <span className="ml-3">{t('使用 微信 继续')}</span>
                  </Button>
                )}

                {status.github_oauth && (
                  <Button
                    theme='outline'
                    className="w-full h-12 flex items-center justify-center !rounded-full border border-gray-200 hover:bg-gray-50 transition-colors"
                    type="tertiary"
                    icon={<IconGithubLogo size="large" />}
                    size="large"
                    onClick={handleGitHubClick}
                    loading={githubLoading}
                  >
                    <span className="ml-3">{t('使用 GitHub 继续')}</span>
                  </Button>
                )}

                {status.oidc_enabled && (
                  <Button
                    theme='outline'
                    className="w-full h-12 flex items-center justify-center !rounded-full border border-gray-200 hover:bg-gray-50 transition-colors"
                    type="tertiary"
                    icon={<OIDCIcon style={{ color: '#1877F2' }} />}
                    size="large"
                    onClick={handleOIDCClick}
                    loading={oidcLoading}
                  >
                    <span className="ml-3">{t('使用 OIDC 继续')}</span>
                  </Button>
                )}

                {status.linuxdo_oauth && (
                  <Button
                    theme='outline'
                    className="w-full h-12 flex items-center justify-center !rounded-full border border-gray-200 hover:bg-gray-50 transition-colors"
                    type="tertiary"
                    icon={<LinuxDoIcon style={{ color: '#E95420', width: '20px', height: '20px' }} />}
                    size="large"
                    onClick={handleLinuxDOClick}
                    loading={linuxdoLoading}
                  >
                    <span className="ml-3">{t('使用 LinuxDO 继续')}</span>
                  </Button>
                )}

                {status.telegram_oauth && (
                  <div className="flex justify-center my-2">
                    <TelegramLoginButton
                      dataOnauth={onTelegramLoginClicked}
                      botName={status.telegram_bot_name}
                    />
                  </div>
                )}

                <Divider margin='12px' align='center'>
                  {t('或')}
                </Divider>

                <Button
                  theme="solid"
                  type="primary"
                  className="w-full h-12 flex items-center justify-center bg-black text-white !rounded-full hover:bg-gray-800 transition-colors"
                  icon={<IconPhone size="large" />}
                  size="large"
                  onClick={handlePhoneRegisterClick}
                  loading={phoneRegisterLoading}
                >
                  <span className="ml-3">{t('使用 手机号 注册')}</span>
                </Button>
              </div>

              <div className="mt-6 text-center text-sm">
                <Text>{t('已有账户？')} <Link to="/login" className="text-blue-600 hover:text-blue-800 font-medium">{t('登录')}</Link></Text>
              </div>
            </div>
          </Card>
        </div>
      </div>
    );
  };

  const renderPhoneRegisterForm = () => {
    return (
      <div className="flex flex-col items-center">
        <div className="w-full max-w-md">
          <div className="flex items-center justify-center mb-6 gap-2">
            <img src={logo} alt="Logo" className="h-10 rounded-full" />
            <Title heading={3} className='!text-gray-800'>{systemName}</Title>
          </div>

          <Card className="shadow-xl border-0 !rounded-2xl overflow-hidden">
            <div className="flex justify-center pt-6 pb-2">
              <Title heading={3} className="text-gray-800 dark:text-gray-200">{t('注 册')}</Title>
            </div>
            <div className="px-2 py-8">
              <Form className="space-y-3">
                <Form.Input
                  field="phone"
                  label={t('手机号')}
                  placeholder={t('请输入11位手机号')}
                  name="phone"
                  size="large"
                  onChange={(value) => handleChange('phone', value)}
                  prefix={<IconPhone />}
                  maxLength={11}
                />

                <Form.Input
                  field="phone_verification_code"
                  label={t('手机验证码')}
                  placeholder={t('输入验证码')}
                  name="phone_verification_code"
                  size="large"
                  onChange={(value) => handleChange('phone_verification_code', value)}
                  prefix={<IconKey />}
                  suffix={
                    <Button
                      onClick={sendVerificationCode}
                      loading={verificationCodeLoading}
                      size="small"
                    >
                      {t('获取验证码')}
                    </Button>
                  }
                />

                <Form.Input
                  field="display_name"
                  label={t('显示名称')}
                  placeholder={t('请输入显示名称')}
                  name="display_name"
                  size="large"
                  onChange={(value) => handleChange('display_name', value)}
                  prefix={<IconUser />}
                />

                <Form.Input
                  field="school"
                  label={t('学校')}
                  placeholder={t('请输入学校名称（可选）')}
                  name="school"
                  size="large"
                  onChange={(value) => handleChange('school', value)}
                  prefix={<IconUser />}
                />

                <Form.Input
                  field="college"
                  label={t('学院')}
                  placeholder={t('请输入学院名称（可选）')}
                  name="college"
                  size="large"
                  onChange={(value) => handleChange('college', value)}
                  prefix={<IconUser />}
                />



                <div className="space-y-2 pt-2">
                  <Button
                    theme="solid"
                    className="w-full !rounded-full"
                    type="primary"
                    htmlType="submit"
                    size="large"
                    onClick={handleSubmit}
                    loading={registerLoading}
                  >
                    {t('注册')}
                  </Button>
                </div>
              </Form>

              {(status.github_oauth || status.oidc_enabled || status.wechat_login || status.linuxdo_oauth || status.telegram_oauth) && (
                <>
                  <Divider margin='12px' align='center'>
                    {t('或')}
                  </Divider>

                  <div className="mt-4 text-center">
                    <Button
                      theme="outline"
                      type="tertiary"
                      className="w-full !rounded-full"
                      size="large"
                      onClick={handleOtherRegisterOptionsClick}
                      loading={otherRegisterOptionsLoading}
                    >
                      {t('其他注册选项')}
                    </Button>
                  </div>
                </>
              )}

              <div className="mt-6 text-center text-sm">
                <Text>{t('已有账户？')} <Link to="/login" className="text-blue-600 hover:text-blue-800 font-medium">{t('登录')}</Link></Text>
              </div>
            </div>
          </Card>
        </div>
      </div>
    );
  };

  const renderWeChatLoginModal = () => {
    return (
      <Modal
        title={t('微信扫码登录')}
        visible={showWeChatLoginModal}
        maskClosable={true}
        onOk={onSubmitWeChatVerificationCode}
        onCancel={() => setShowWeChatLoginModal(false)}
        okText={t('登录')}
        size="small"
        centered={true}
        okButtonProps={{
          loading: wechatCodeSubmitLoading,
        }}
      >
        <div className="flex flex-col items-center">
          <img src={status.wechat_qrcode} alt="微信二维码" className="mb-4" />
        </div>

        <div className="text-center mb-4">
          <p>{t('微信扫码关注公众号，输入「验证码」获取验证码（三分钟内有效）')}</p>
        </div>

        <Form size="large">
          <Form.Input
            field="wechat_verification_code"
            placeholder={t('验证码')}
            label={t('验证码')}
            value={inputs.wechat_verification_code}
            onChange={(value) => handleChange('wechat_verification_code', value)}
          />
        </Form>
      </Modal>
    );
  };

  return (
    <div className="relative overflow-hidden bg-gray-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      {/* 背景模糊晕染球 */}
      <div className="blur-ball blur-ball-indigo" style={{ top: '-80px', right: '-80px', transform: 'none' }} />
      <div className="blur-ball blur-ball-teal" style={{ top: '50%', left: '-120px' }} />
      <div className="w-full max-w-sm mt-[64px]">
        {showPhoneRegister || !(status.github_oauth || status.oidc_enabled || status.wechat_login || status.linuxdo_oauth || status.telegram_oauth)
          ? renderPhoneRegisterForm()
          : renderOAuthOptions()}
        {renderWeChatLoginModal()}

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
  );
};

export default RegisterForm;
