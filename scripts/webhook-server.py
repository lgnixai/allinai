#!/usr/bin/env python3
"""
Webhook服务器
用于接收GitHub/GitLab的webhook通知并触发自动部署
"""

import json
import subprocess
import hmac
import hashlib
import os
from flask import Flask, request, jsonify
from datetime import datetime

app = Flask(__name__)

# 配置
WEBHOOK_SECRET = os.getenv('WEBHOOK_SECRET', 'your_webhook_secret')  # 设置webhook密钥
PROJECT_DIR = os.getenv('PROJECT_DIR', '/path/to/your/project')  # 项目目录
DEPLOY_SCRIPT = os.path.join(PROJECT_DIR, 'scripts/deploy.sh')
LOG_FILE = '/var/log/webhook.log'

def log(message):
    """记录日志"""
    timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    log_entry = f"[{timestamp}] {message}\n"
    with open(LOG_FILE, 'a') as f:
        f.write(log_entry)
    print(log_entry.strip())

def verify_signature(payload, signature):
    """验证webhook签名"""
    if not signature:
        return False
    
    # GitHub格式: sha256=xxx
    if signature.startswith('sha256='):
        signature = signature[7:]
    
    expected_signature = hmac.new(
        WEBHOOK_SECRET.encode('utf-8'),
        payload,
        hashlib.sha256
    ).hexdigest()
    
    return hmac.compare_digest(signature, expected_signature)

@app.route('/webhook', methods=['POST'])
def webhook():
    """处理webhook请求"""
    try:
        # 获取请求数据
        payload = request.get_data()
        signature = request.headers.get('X-Hub-Signature-256') or request.headers.get('X-Hub-Signature')
        
        # 验证签名
        if not verify_signature(payload, signature):
            log("签名验证失败")
            return jsonify({'error': 'Invalid signature'}), 401
        
        # 解析JSON数据
        data = json.loads(payload)
        event_type = request.headers.get('X-GitHub-Event') or request.headers.get('X-Gitlab-Event')
        
        log(f"收到 {event_type} 事件")
        
        # 处理push事件
        if event_type == 'push':
            ref = data.get('ref', '')
            branch = ref.replace('refs/heads/', '')
            
            log(f"检测到分支推送: {branch}")
            
            # 只对main分支进行自动部署
            if branch == 'main':
                log("开始自动部署...")
                
                # 执行部署脚本
                if os.path.exists(DEPLOY_SCRIPT):
                    result = subprocess.run(
                        [DEPLOY_SCRIPT, branch],
                        cwd=PROJECT_DIR,
                        capture_output=True,
                        text=True
                    )
                    
                    if result.returncode == 0:
                        log("部署成功")
                        return jsonify({'status': 'success', 'message': 'Deployment completed'})
                    else:
                        log(f"部署失败: {result.stderr}")
                        return jsonify({'status': 'error', 'message': result.stderr}), 500
                else:
                    log(f"部署脚本不存在: {DEPLOY_SCRIPT}")
                    return jsonify({'status': 'error', 'message': 'Deploy script not found'}), 500
            else:
                log(f"跳过分支 {branch} 的自动部署")
                return jsonify({'status': 'skipped', 'message': f'Skipped branch {branch}'})
        
        return jsonify({'status': 'ignored', 'message': f'Ignored event type: {event_type}'})
        
    except Exception as e:
        log(f"处理webhook时出错: {str(e)}")
        return jsonify({'error': str(e)}), 500

@app.route('/health', methods=['GET'])
def health():
    """健康检查"""
    return jsonify({'status': 'healthy', 'timestamp': datetime.now().isoformat()})

if __name__ == '__main__':
    log("Webhook服务器启动")
    app.run(host='0.0.0.0', port=8080, debug=False)

