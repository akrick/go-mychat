-- 咨询师评价表
CREATE TABLE IF NOT EXISTS counselor_reviews (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id INT UNSIGNED NOT NULL UNIQUE,
    user_id INT UNSIGNED NOT NULL,
    counselor_id INT UNSIGNED NOT NULL,
    rating INT NOT NULL COMMENT '评分(1-5)',
    comment TEXT COMMENT '评价内容',
    is_anonymous BOOLEAN DEFAULT FALSE COMMENT '是否匿名',
    is_visible BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order_id (order_id),
    INDEX idx_user_id (user_id),
    INDEX idx_counselor_id (counselor_id),
    INDEX idx_is_visible (is_visible)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='咨询师评价表';

-- 插入测试评价数据
INSERT INTO counselor_reviews (order_id, user_id, counselor_id, rating, comment, is_anonymous, is_visible) VALUES
-- 咨询师1的评价
(1, 1, 1, 5, '王老师非常专业，耐心倾听我的问题，给出了很多有用的建议。通过几次咨询，我对自己的情况有了更清晰的认识，非常感谢！', FALSE, TRUE),
(2, 2, 1, 5, '咨询体验很好，王老师能够准确抓住问题要点，引导我思考，帮助我找到了解决问题的方向。', FALSE, TRUE),
(3, 3, 1, 4, '第一次咨询感觉不错，王老师很有亲和力，让我能够敞开心扉。希望后续的咨询能够有更多实质性的帮助。', FALSE, TRUE),
(4, 4, 1, 5, '非常感谢王老师的帮助，她的专业素养和人文关怀让我深受感动。强烈推荐！', TRUE, TRUE),
(5, 5, 1, 5, '咨询效果超出预期，王老师的分析方法很独到，帮助我从不同的角度看待问题。', FALSE, TRUE),

-- 咨询师2的评价
(6, 6, 2, 5, '李博士的咨询非常有深度，他的理论功底很扎实，能够将复杂的心理问题用简单易懂的方式解释清楚。', FALSE, TRUE),
(7, 7, 2, 4, '李博士很专业，但有时候表达方式比较学术化，需要多交流才能理解。整体还是很满意的。', FALSE, TRUE),
(8, 8, 2, 5, '李博士的洞察力很强，能够快速找到问题的根源。咨询后我对自己有了更深入的认识。', FALSE, TRUE),

-- 咨询师3的评价
(9, 9, 3, 5, '张老师的咨询风格很温暖，让我感受到了真正的关怀。她不仅帮助我解决了问题，还教会了我很多自我调节的方法。', FALSE, TRUE),
(10, 10, 3, 5, '张老师非常有耐心，每次咨询都会认真听我说话，给我充分的表达空间。非常推荐！', FALSE, TRUE),
(11, 11, 3, 4, '咨询效果不错，张老师很专业。希望能够在后续的咨询中更多地学习一些实用的心理技巧。', FALSE, TRUE),

-- 咨询师4的评价
(12, 12, 4, 5, '刘老师很有经验，对职场心理问题特别了解。经过几次咨询，我明显感觉自己的心态变得更加积极了。', FALSE, TRUE),
(13, 13, 4, 4, '刘老师的建议很实用，但我感觉还需要更多的咨询时间来解决深层次的问题。整体还是很满意的。', FALSE, TRUE),

-- 咨询师5的评价
(14, 14, 5, 5, '陈老师是我遇到过的最好的咨询师，她不仅专业，而且非常真诚。通过咨询，我找回了对生活的信心。', FALSE, TRUE),
(15, 15, 5, 5, '陈老师的方法很有效，她不会直接告诉我该怎么做，而是引导我自己找到答案。这种方式让我受益匪浅。', FALSE, TRUE),
(16, 16, 5, 5, '咨询体验非常好，陈老师能够准确理解我的感受，给予我支持和鼓励。', TRUE, TRUE),

-- 咨询师6的评价
(17, 17, 6, 4, '杨老师的咨询很专业，她的分析方法很系统。有时候感觉时间有点短，希望能够有更深入的交流。', FALSE, TRUE),
(18, 18, 6, 5, '杨老师很有亲和力，让我很快就能放松下来。咨询效果很好，问题得到了很好的解决。', FALSE, TRUE),

-- 咨询师7的评价
(19, 19, 7, 5, '周老师特别擅长家庭心理问题，她帮助我和家人重建了良好的沟通方式。非常感谢！', FALSE, TRUE),
(20, 20, 7, 5, '周老师的咨询很温暖，她不仅解决了我的问题，还让我学会了如何更好地和家人相处。', FALSE, TRUE),

-- 咨询师8的评价
(21, 21, 8, 5, '吴老师是专业的婚姻咨询师，她帮助我和伴侣化解了很多矛盾。现在我们的关系比以前更好了。', FALSE, TRUE),
(22, 22, 8, 4, '吴老师的建议很有建设性，她的婚姻理论很实用。虽然咨询时间不长，但已经看到了积极的变化。', FALSE, TRUE);
