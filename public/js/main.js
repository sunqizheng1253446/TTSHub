// TTSHub 控制台主脚本

// DOM 加载完成后初始化
 document.addEventListener('DOMContentLoaded', () => {
    // 初始化应用
    initApp();
});

/**
 * 初始化应用
 */
 function initApp() {
    // 初始化主题
    initTheme();
    
    // 初始化导航
    initNavigation();
    
    // 初始化模态框
    initModal();
    
    // 初始化测试面板
    initTestPanel();
    
    // 初始化通知系统
    initNotification();
    
    // 初始化系统设置
    initSettings();
    
    // 加载渠道列表
    loadChannels();
    
    // 加载统计数据
    loadDashboardStats();
    
    // 加载最近活动
    loadRecentActivity();
    
    // 加载日志数据
    loadLogs();
}

/**
 * 初始化主题
 */
 function initTheme() {
    const themeBtn = document.getElementById('theme-toggle');
    
    // 从本地存储加载主题设置
    const savedTheme = localStorage.getItem('ttshub-theme');
    if (savedTheme === 'dark') {
        document.body.classList.add('dark-theme');
    }
    
    // 主题切换事件
    themeBtn.addEventListener('click', () => {
        document.body.classList.toggle('dark-theme');
        const currentTheme = document.body.classList.contains('dark-theme') ? 'dark' : 'light';
        localStorage.setItem('ttshub-theme', currentTheme);
        
        // 更新图标
        updateThemeIcon();
    });
    
    // 初始更新图标
    updateThemeIcon();
}

/**
 * 更新主题图标
 */
 function updateThemeIcon() {
    const themeIcon = document.getElementById('theme-icon');
    if (document.body.classList.contains('dark-theme')) {
        themeIcon.textContent = '☀️';
    } else {
        themeIcon.textContent = '🌙';
    }
}

/**
 * 初始化导航
 */
 function initNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    const panels = document.querySelectorAll('.panel');
    
    // 导航点击事件
    navItems.forEach(item => {
        item.addEventListener('click', () => {
            const panelId = item.getAttribute('data-panel');
            
            // 更新活动状态
            navItems.forEach(nav => nav.classList.remove('active'));
            item.classList.add('active');
            
            // 显示对应面板
            panels.forEach(panel => {
                panel.classList.remove('active');
                if (panel.id === panelId) {
                    panel.classList.add('active');
                }
            });
            
            // 保存当前面板到本地存储
            localStorage.setItem('ttshub-current-panel', panelId);
        });
    });
    
    // 从本地存储恢复最后访问的面板
    const savedPanel = localStorage.getItem('ttshub-current-panel');
    if (savedPanel) {
        const savedNavItem = document.querySelector(`.nav-item[data-panel="${savedPanel}"]`);
        if (savedNavItem) {
            savedNavItem.click();
        } else {
            // 如果保存的面板不存在，默认显示第一个面板
            navItems[0].click();
        }
    } else {
        // 默认显示第一个面板
        navItems[0].click();
    }
}

/**
 * 初始化模态框
 */
 function initModal() {
    const modal = document.getElementById('channel-modal');
    const modalClose = document.getElementById('modal-close');
    const modalBackdrop = document.querySelector('.modal-backdrop');
    const addChannelBtn = document.getElementById('add-channel-btn');
    const saveChannelBtn = document.getElementById('save-channel-btn');
    const cancelChannelBtn = document.getElementById('cancel-channel-btn');
    
    // 打开模态框
    const openModal = (mode, channelData = null) => {
        modal.classList.add('active');
        
        // 设置模态框标题
        const modalTitle = document.getElementById('modal-title');
        const channelIdField = document.getElementById('channel-id');
        
        if (mode === 'add') {
            modalTitle.textContent = '添加新渠道';
            channelIdField.value = '';
            channelIdField.disabled = false;
            // 重置表单
            document.getElementById('channel-form').reset();
        } else {
            modalTitle.textContent = '编辑渠道';
            // 填充表单数据
            if (channelData) {
                channelIdField.value = channelData.id;
                channelIdField.disabled = true;
                document.getElementById('channel-name').value = channelData.name || '';
                document.getElementById('channel-type').value = channelData.type || '';
                document.getElementById('channel-status').value = channelData.status || 'inactive';
                document.getElementById('channel-api-key').value = channelData.apiKey || '';
                document.getElementById('channel-api-secret').value = channelData.apiSecret || '';
                document.getElementById('channel-endpoint').value = channelData.endpoint || '';
                document.getElementById('channel-model').value = channelData.model || '';
                document.getElementById('channel-voice').value = channelData.voice || '';
                document.getElementById('channel-speed').value = channelData.speed || '1.0';
                document.getElementById('channel-pitch').value = channelData.pitch || '1.0';
            }
        }
        
        // 设置保存按钮的模式
        saveChannelBtn.setAttribute('data-mode', mode);
    };
    
    // 关闭模态框
    const closeModal = () => {
        modal.classList.remove('active');
    };
    
    // 添加渠道按钮点击事件
    addChannelBtn.addEventListener('click', () => openModal('add'));
    
    // 关闭按钮点击事件
    modalClose.addEventListener('click', closeModal);
    
    // 点击背景关闭
    modalBackdrop.addEventListener('click', closeModal);
    
    // 取消按钮点击事件
    cancelChannelBtn.addEventListener('click', closeModal);
    
    // 保存按钮点击事件
    saveChannelBtn.addEventListener('click', () => {
        const mode = saveChannelBtn.getAttribute('data-mode');
        const channelData = collectChannelFormData();
        
        if (mode === 'add') {
            addChannel(channelData);
        } else {
            updateChannel(channelData);
        }
    });
    
    // 处理渠道操作按钮点击
    document.addEventListener('click', (e) => {
        if (e.target.closest('.edit-channel-btn')) {
            const channelId = e.target.closest('.channel-card').dataset.id;
            const channelData = getChannelById(channelId);
            openModal('edit', channelData);
        } else if (e.target.closest('.delete-channel-btn')) {
            const channelId = e.target.closest('.channel-card').dataset.id;
            deleteChannel(channelId);
        } else if (e.target.closest('.test-channel-btn')) {
            const channelId = e.target.closest('.channel-card').dataset.id;
            testChannel(channelId);
        }
    });
    
    // 导出打开模态框函数供其他地方使用
    window.openModal = openModal;
}

/**
 * 收集渠道表单数据
 */
 function collectChannelFormData() {
    return {
        id: document.getElementById('channel-id').value,
        name: document.getElementById('channel-name').value,
        type: document.getElementById('channel-type').value,
        status: document.getElementById('channel-status').value,
        apiKey: document.getElementById('channel-api-key').value,
        apiSecret: document.getElementById('channel-api-secret').value,
        endpoint: document.getElementById('channel-endpoint').value,
        model: document.getElementById('channel-model').value,
        voice: document.getElementById('channel-voice').value,
        speed: document.getElementById('channel-speed').value,
        pitch: document.getElementById('channel-pitch').value
    };
}

/**
 * 初始化测试面板
 */
 function initTestPanel() {
    const testForm = document.getElementById('tts-test-form');
    const textInput = document.getElementById('test-text');
    const charCount = document.getElementById('char-count');
    const synthesizeBtn = document.getElementById('synthesize-btn');
    const clearBtn = document.getElementById('clear-test-btn');
    
    // 字符计数
    textInput.addEventListener('input', () => {
        charCount.textContent = `${textInput.value.length} 字符`;
    });
    
    // 合成按钮点击事件
    synthesizeBtn.addEventListener('click', (e) => {
        e.preventDefault();
        synthesizeText();
    });
    
    // 清空按钮点击事件
    clearBtn.addEventListener('click', () => {
        testForm.reset();
        charCount.textContent = '0 字符';
        document.getElementById('audio-result').innerHTML = '<div class="empty-player">暂无测试结果</div>';
        document.getElementById('test-details').innerHTML = '';
    });
    
    // 语速和音高滑块
    const speedSlider = document.getElementById('test-speed');
    const pitchSlider = document.getElementById('test-pitch');
    const speedValue = document.getElementById('speed-value');
    const pitchValue = document.getElementById('pitch-value');
    
    speedSlider.addEventListener('input', () => {
        speedValue.textContent = speedSlider.value;
    });
    
    pitchSlider.addEventListener('input', () => {
        pitchValue.textContent = pitchSlider.value;
    });
}

/**
 * 初始化通知系统
 */
 function initNotification() {
    // 导出通知函数
    window.showNotification = showNotification;
}

/**
 * 显示通知
 */
 function showNotification(type, message) {
    const notification = document.getElementById('notification');
    const notificationMessage = document.getElementById('notification-message');
    
    // 设置通知类型和消息
    notification.className = `notification ${type} active`;
    notificationMessage.textContent = message;
    
    // 3秒后自动关闭
    setTimeout(() => {
        notification.classList.remove('active');
    }, 3000);
}

/**
 * 初始化系统设置
 */
 function initSettings() {
    const settingTabs = document.querySelectorAll('.tab');
    const tabContents = document.querySelectorAll('.tab-content');
    
    // 设置标签点击事件
    settingTabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const tabId = tab.getAttribute('data-tab');
            
            // 更新活动状态
            settingTabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');
            
            // 显示对应内容
            tabContents.forEach(content => {
                content.classList.remove('active');
                if (content.id === tabId) {
                    content.classList.add('active');
                }
            });
        });
    });
    
    // 默认激活第一个标签
    settingTabs[0].click();
    
    // 数据库操作按钮
    document.getElementById('backup-db-btn').addEventListener('click', backupDatabase);
    document.getElementById('restore-db-btn').addEventListener('click', restoreDatabase);
    document.getElementById('reset-db-btn').addEventListener('click', resetDatabase);
}

/**
 * 加载渠道列表
 */
 function loadChannels() {
    // 显示加载中
    const channelsList = document.getElementById('channels-list');
    channelsList.innerHTML = '<div class="loading-spinner"></div>';
    
    // 模拟API调用
    setTimeout(() => {
        // 从本地存储获取渠道数据（实际项目中应从API获取）
        const channels = getChannelsFromStorage();
        
        if (channels.length === 0) {
            // 显示空状态
            channelsList.innerHTML = `
                <div class="empty-state">
                    <div class="empty-icon">📢</div>
                    <p>暂无渠道配置</p>
                    <p class="empty-hint">点击上方"添加渠道"按钮创建您的第一个TTS渠道</p>
                </div>
            `;
        } else {
            // 渲染渠道列表
            channelsList.innerHTML = channels.map(channel => `
                <div class="channel-card" data-id="${channel.id}">
                    <div class="channel-card-header">
                        <div class="channel-info">
                            <h4>${channel.name}</h4>
                            <span class="channel-type">${channel.type}</span>
                        </div>
                        <div class="channel-actions">
                            <button class="channel-action-btn test-channel-btn" title="测试渠道">🔊</button>
                            <button class="channel-action-btn edit-channel-btn" title="编辑渠道">✏️</button>
                            <button class="channel-action-btn delete-channel-btn" title="删除渠道">🗑️</button>
                        </div>
                    </div>
                    <div class="channel-details">
                        <div class="channel-detail-item">
                            <span>状态</span>
                            <span class="channel-status ${channel.status}">${channel.status === 'active' ? '活跃' : '未激活'}</span>
                        </div>
                        <div class="channel-detail-item">
                            <span>模型</span>
                            <span>${channel.model || '默认'}</span>
                        </div>
                        <div class="channel-detail-item">
                            <span>语音</span>
                            <span>${channel.voice || '默认'}</span>
                        </div>
                        <div class="channel-detail-item">
                            <span>创建时间</span>
                            <span>${channel.createdAt || '未知'}</span>
                        </div>
                    </div>
                </div>
            `).join('');
        }
    }, 800);
}

/**
 * 从本地存储获取渠道列表
 */
 function getChannelsFromStorage() {
    const channels = localStorage.getItem('ttshub-channels');
    return channels ? JSON.parse(channels) : [];
}

/**
 * 保存渠道列表到本地存储
 */
 function saveChannelsToStorage(channels) {
    localStorage.setItem('ttshub-channels', JSON.stringify(channels));
}

/**
 * 根据ID获取渠道
 */
 function getChannelById(channelId) {
    const channels = getChannelsFromStorage();
    return channels.find(channel => channel.id === channelId);
}

/**
 * 添加渠道
 */
 function addChannel(channelData) {
    showLoading(true);
    
    // 模拟API调用
    setTimeout(() => {
        try {
            const channels = getChannelsFromStorage();
            
            // 检查ID是否已存在
            if (channels.some(c => c.id === channelData.id)) {
                showNotification('error', '渠道ID已存在');
                return;
            }
            
            // 添加创建时间
            channelData.createdAt = new Date().toISOString().slice(0, 10);
            
            // 添加到列表
            channels.push(channelData);
            saveChannelsToStorage(channels);
            
            // 刷新列表
            loadChannels();
            
            // 关闭模态框
            document.getElementById('channel-modal').classList.remove('active');
            
            showNotification('success', '渠道添加成功');
        } catch (error) {
            showNotification('error', '渠道添加失败：' + error.message);
        } finally {
            showLoading(false);
        }
    }, 1000);
}

/**
 * 更新渠道
 */
 function updateChannel(channelData) {
    showLoading(true);
    
    // 模拟API调用
    setTimeout(() => {
        try {
            const channels = getChannelsFromStorage();
            const index = channels.findIndex(c => c.id === channelData.id);
            
            if (index === -1) {
                showNotification('error', '渠道不存在');
                return;
            }
            
            // 保留创建时间
            channelData.createdAt = channels[index].createdAt;
            
            // 更新渠道
            channels[index] = channelData;
            saveChannelsToStorage(channels);
            
            // 刷新列表
            loadChannels();
            
            // 关闭模态框
            document.getElementById('channel-modal').classList.remove('active');
            
            showNotification('success', '渠道更新成功');
        } catch (error) {
            showNotification('error', '渠道更新失败：' + error.message);
        } finally {
            showLoading(false);
        }
    }, 1000);
}

/**
 * 删除渠道
 */
 function deleteChannel(channelId) {
    if (confirm('确定要删除这个渠道吗？此操作不可撤销。')) {
        showLoading(true);
        
        // 模拟API调用
        setTimeout(() => {
            try {
                const channels = getChannelsFromStorage();
                const filteredChannels = channels.filter(c => c.id !== channelId);
                
                if (channels.length === filteredChannels.length) {
                    showNotification('error', '渠道不存在');
                    return;
                }
                
                saveChannelsToStorage(filteredChannels);
                
                // 刷新列表
                loadChannels();
                
                showNotification('success', '渠道删除成功');
            } catch (error) {
                showNotification('error', '渠道删除失败：' + error.message);
            } finally {
                showLoading(false);
            }
        }, 800);
    }
}

/**
 * 测试渠道
 */
 function testChannel(channelId) {
    showLoading(true);
    
    // 模拟API调用
    setTimeout(() => {
        try {
            const channel = getChannelById(channelId);
            
            if (!channel) {
                showNotification('error', '渠道不存在');
                return;
            }
            
            if (channel.status !== 'active') {
                showNotification('warning', '请先激活渠道后再进行测试');
                return;
            }
            
            showNotification('info', `渠道 ${channel.name} 测试成功`);
        } catch (error) {
            showNotification('error', '渠道测试失败：' + error.message);
        } finally {
            showLoading(false);
        }
    }, 1200);
}

/**
 * 合成文本
 */
 function synthesizeText() {
    showLoading(true);
    
    // 收集表单数据
    const channelId = document.getElementById('test-channel').value;
    const text = document.getElementById('test-text').value;
    const voice = document.getElementById('test-voice').value;
    const speed = document.getElementById('test-speed').value;
    const pitch = document.getElementById('test-pitch').value;
    
    // 验证输入
    if (!channelId) {
        showNotification('error', '请选择一个渠道');
        showLoading(false);
        return;
    }
    
    if (!text.trim()) {
        showNotification('error', '请输入要转换的文本');
        showLoading(false);
        return;
    }
    
    // 模拟API调用
    setTimeout(() => {
        try {
            const channel = getChannelById(channelId);
            
            if (!channel) {
                showNotification('error', '选择的渠道不存在');
                return;
            }
            
            if (channel.status !== 'active') {
                showNotification('warning', '选择的渠道未激活');
                return;
            }
            
            // 模拟生成音频URL（实际项目中应返回真实的音频URL）
            const audioUrl = `data:audio/wav;base64,${btoa('dummy audio data')}`;
            
            // 显示结果
            const audioResult = document.getElementById('audio-result');
            audioResult.innerHTML = `
                <div class="audio-player">
                    <audio controls>
                        <source src="${audioUrl}" type="audio/wav">
                        您的浏览器不支持音频元素。
                    </audio>
                </div>
            `;
            
            // 显示详细信息
            const testDetails = document.getElementById('test-details');
            testDetails.innerHTML = `
                <div class="test-info">
                    <div class="info-item">
                        <span class="info-label">使用渠道</span>
                        <span>${channel.name} (${channel.type})</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">语音</span>
                        <span>${voice || '默认'}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">语速</span>
                        <span>${speed}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">音高</span>
                        <span>${pitch}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">文本长度</span>
                        <span>${text.length} 字符</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">响应时间</span>
                        <span>1.2 秒</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">生成时间</span>
                        <span>${new Date().toLocaleString()}</span>
                    </div>
                </div>
            `;
            
            // 记录最近活动
            recordActivity(`使用 ${channel.name} 渠道合成了 ${text.length} 字符的文本`);
            
            // 刷新统计数据
            updateSynthesisCount();
            
            showNotification('success', '文本合成成功');
        } catch (error) {
            showNotification('error', '文本合成失败：' + error.message);
        } finally {
            showLoading(false);
        }
    }, 2000);
}

/**
 * 加载仪表盘统计数据
 */
 function loadDashboardStats() {
    // 模拟API调用
    setTimeout(() => {
        const stats = {
            totalChannels: localStorage.getItem('ttshub-total-channels') || 0,
            activeChannels: localStorage.getItem('ttshub-active-channels') || 0,
            totalSyntheses: localStorage.getItem('ttshub-total-syntheses') || 0,
            todaySyntheses: localStorage.getItem('ttshub-today-syntheses') || 0
        };
        
        // 更新统计卡片
        document.getElementById('total-channels').textContent = stats.totalChannels;
        document.getElementById('active-channels').textContent = stats.activeChannels;
        document.getElementById('total-syntheses').textContent = stats.totalSyntheses;
        document.getElementById('today-syntheses').textContent = stats.todaySyntheses;
    }, 500);
}

/**
 * 更新合成次数统计
 */
 function updateSynthesisCount() {
    // 更新总合成次数
    const totalSyntheses = parseInt(localStorage.getItem('ttshub-total-syntheses') || '0') + 1;
    localStorage.setItem('ttshub-total-syntheses', totalSyntheses.toString());
    document.getElementById('total-syntheses').textContent = totalSyntheses;
    
    // 更新今日合成次数
    const today = new Date().toISOString().slice(0, 10);
    const lastDate = localStorage.getItem('ttshub-last-synthesis-date');
    let todaySyntheses = 0;
    
    if (lastDate === today) {
        todaySyntheses = parseInt(localStorage.getItem('ttshub-today-syntheses') || '0') + 1;
    } else {
        todaySyntheses = 1;
        localStorage.setItem('ttshub-last-synthesis-date', today);
    }
    
    localStorage.setItem('ttshub-today-syntheses', todaySyntheses.toString());
    document.getElementById('today-syntheses').textContent = todaySyntheses;
}

/**
 * 加载最近活动
 */
 function loadRecentActivity() {
    const activityList = document.getElementById('activity-list');
    
    // 从本地存储获取活动记录
    const activities = getActivitiesFromStorage();
    
    if (activities.length === 0) {
        activityList.innerHTML = '<div class="empty-state">暂无活动记录</div>';
    } else {
        activityList.innerHTML = activities.map(activity => `
            <div class="activity-item">
                <div class="activity-time">${formatTime(activity.time)}</div>
                <div class="activity-content">${activity.message}</div>
            </div>
        `).join('');
    }
}

/**
 * 从本地存储获取活动记录
 */
 function getActivitiesFromStorage() {
    const activities = localStorage.getItem('ttshub-activities');
    return activities ? JSON.parse(activities) : [];
}

/**
 * 保存活动记录到本地存储
 */
 function saveActivitiesToStorage(activities) {
    // 最多保存50条记录
    if (activities.length > 50) {
        activities = activities.slice(0, 50);
    }
    localStorage.setItem('ttshub-activities', JSON.stringify(activities));
}

/**
 * 记录活动
 */
 function recordActivity(message) {
    const activities = getActivitiesFromStorage();
    activities.unshift({
        time: new Date().toISOString(),
        message: message
    });
    saveActivitiesToStorage(activities);
    loadRecentActivity();
}

/**
 * 格式化时间
 */
 function formatTime(timeString) {
    const date = new Date(timeString);
    const now = new Date();
    const diffInMinutes = Math.floor((now - date) / (1000 * 60));
    
    if (diffInMinutes < 1) return '刚刚';
    if (diffInMinutes < 60) return `${diffInMinutes}分钟前`;
    if (diffInMinutes < 24 * 60) return `${Math.floor(diffInMinutes / 60)}小时前`;
    
    return date.toLocaleDateString();
}

/**
 * 加载日志
 */
 function loadLogs() {
    const logsList = document.getElementById('logs-list');
    
    // 模拟API调用
    setTimeout(() => {
        const logs = [
            {
                time: new Date().toISOString().slice(0, 19).replace('T', ' '),
                level: 'info',
                message: 'TTSHub服务启动成功'
            },
            {
                time: new Date(Date.now() - 3600000).toISOString().slice(0, 19).replace('T', ' '),
                level: 'debug',
                message: '加载配置文件: config.yaml'
            },
            {
                time: new Date(Date.now() - 7200000).toISOString().slice(0, 19).replace('T', ' '),
                level: 'warn',
                message: '部分渠道API密钥即将过期'
            },
            {
                time: new Date(Date.now() - 86400000).toISOString().slice(0, 19).replace('T', ' '),
                level: 'info',
                message: '数据库备份完成'
            },
            {
                time: new Date(Date.now() - 172800000).toISOString().slice(0, 19).replace('T', ' '),
                level: 'error',
                message: '渠道请求失败: API密钥无效'
            }
        ];
        
        logsList.innerHTML = logs.map(log => `
            <div class="log-item">
                <div class="log-time">${log.time}</div>
                <div class="log-content">
                    <span class="log-level ${log.level}">${log.level.toUpperCase()}</span>
                    ${log.message}
                </div>
            </div>
        `).join('');
    }, 800);
}

/**
 * 数据库备份
 */
 function backupDatabase() {
    showLoading(true);
    
    // 模拟API调用
    setTimeout(() => {
        try {
            // 在实际项目中，这里应该调用后端API进行数据库备份
            const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, 19);
            const filename = `ttshub-backup-${timestamp}.sql`;
            
            showNotification('success', `数据库备份成功，文件名：${filename}`);
            recordActivity(`数据库备份完成：${filename}`);
        } catch (error) {
            showNotification('error', '数据库备份失败：' + error.message);
        } finally {
            showLoading(false);
        }
    }, 2000);
}

/**
 * 数据库恢复
 */
 function restoreDatabase() {
    if (confirm('确定要恢复数据库吗？这将覆盖当前所有数据，操作不可撤销！')) {
        // 提示用户选择备份文件
        alert('请在实际应用中实现文件选择功能');
        // 在实际项目中，这里应该打开文件选择对话框并调用后端API进行数据库恢复
    }
}

/**
 * 重置数据库
 */
 function resetDatabase() {
    if (confirm('确定要重置数据库吗？这将删除所有配置数据，操作不可撤销！')) {
        showLoading(true);
        
        // 模拟API调用
        setTimeout(() => {
            try {
                // 清除本地存储中的数据
                localStorage.removeItem('ttshub-channels');
                localStorage.removeItem('ttshub-activities');
                localStorage.setItem('ttshub-total-channels', '0');
                localStorage.setItem('ttshub-active-channels', '0');
                localStorage.setItem('ttshub-total-syntheses', '0');
                localStorage.setItem('ttshub-today-syntheses', '0');
                
                // 刷新UI
                loadChannels();
                loadDashboardStats();
                loadRecentActivity();
                
                showNotification('success', '数据库重置成功');
                recordActivity('数据库已重置');
            } catch (error) {
                showNotification('error', '数据库重置失败：' + error.message);
            } finally {
                showLoading(false);
            }
        }, 1500);
    }
}

/**
 * 显示/隐藏加载遮罩
 */
 function showLoading(show) {
    const loadingOverlay = document.getElementById('loading-overlay');
    if (show) {
        loadingOverlay.classList.add('active');
    } else {
        loadingOverlay.classList.remove('active');
    }
}

/**
 * 健康检查
 */
 function checkHealth() {
    // 模拟API调用
    return true;
}

/**
 * 初始化模拟数据（仅在开发环境使用）
 */
 function initMockData() {
    // 检查是否已有数据
    if (!localStorage.getItem('ttshub-channels')) {
        const mockChannels = [
            {
                id: 'openai',
                name: 'OpenAI TTS',
                type: 'openai',
                status: 'active',
                apiKey: 'sk-...1234',
                endpoint: 'https://api.openai.com/v1/audio/speech',
                model: 'tts-1',
                voice: 'alloy',
                speed: '1.0',
                pitch: '1.0',
                createdAt: new Date(Date.now() - 86400000).toISOString().slice(0, 10)
            },
            {
                id: 'baidu',
                name: '百度语音合成',
                type: 'baidu',
                status: 'active',
                apiKey: 'xxx',
                apiSecret: 'yyy',
                endpoint: 'https://tsn.baidu.com/text2audio',
                model: '1',
                voice: '度小宇',
                speed: '5',
                pitch: '5',
                createdAt: new Date(Date.now() - 172800000).toISOString().slice(0, 10)
            },
            {
                id: 'google',
                name: 'Google Cloud TTS',
                type: 'google',
                status: 'inactive',
                apiKey: '',
                endpoint: 'https://texttospeech.googleapis.com/v1/text:synthesize',
                model: 'en-US-Wavenet-A',
                voice: 'en-US-Wavenet-A',
                speed: '1.0',
                pitch: '0.0',
                createdAt: new Date(Date.now() - 259200000).toISOString().slice(0, 10)
            }
        ];
        
        localStorage.setItem('ttshub-channels', JSON.stringify(mockChannels));
        localStorage.setItem('ttshub-total-channels', '3');
        localStorage.setItem('ttshub-active-channels', '2');
        localStorage.setItem('ttshub-total-syntheses', '42');
        localStorage.setItem('ttshub-today-syntheses', '8');
        
        // 初始化活动记录
        const mockActivities = [
            {
                time: new Date().toISOString(),
                message: '系统启动成功'
            },
            {
                time: new Date(Date.now() - 3600000).toISOString(),
                message: '添加了Google Cloud TTS渠道'
            },
            {
                time: new Date(Date.now() - 7200000).toISOString(),
                message: '使用百度语音合成渠道测试成功'
            }
        ];
        
        localStorage.setItem('ttshub-activities', JSON.stringify(mockActivities));
    }
}

// 初始化模拟数据
initMockData();