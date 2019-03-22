import tensorflow as tf
import os
import csv
import numpy
import sys
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '2'

cols = ['ho_score', 'ho_half_score', 'gu_score', 'gu_half_score', 'half_let', 'half_let_hm', 'half_let_aw', 'half_avg_hm', 'half_avg_aw', 'half_avg_eq', 'half_size', 'half_size_big', 'half_size_sma', 'half_first_let', 'half_first_let_hm', 'half_first_let_aw', 'half_first_avg_hm', 'half_first_avg_aw', 'half_first_avg_eq', 'half_first_size', 'half_first_size_big', 'half_first_size_sma', 'let', 'let_hm', 'let_aw', 'avg_hm', 'avg_aw', 'avg_eq', 'size', 'size_big', 'size_sma', 'first_let', 'first_let_hm', 'first_let_aw', 'first_avg_hm', 'first_avg_aw', 'first_avg_eq', 'first_size', 'first_size_big', 'first_size_sma']
class DataSet(object):
    def __init__(self, file):
        with open(file) as f:
            reader = csv.reader(f)
            self.datas = []
            self.labels = []
            self.nums = 0
            self._index_in_epoch = 0
            self._epochs_completed = 0
            self._headers = []
            self._header_keys = {}
            for row in reader:
                if reader.line_num == 1:
                    self._headers = row
                    for i in range(len(row)):
                        self._header_keys[row[i]] = i
                    print(self._headers)
                    continue
                self.preparehalf(row)
                    
            print("load data:" + str(self.nums))

    def getDataLen(self):
        return len(self.datas[0])

    def preparehalf(self, row):
        data = []
        for i in range(len(row)):
            data.append(numpy.float32(row[i]))

        ss1 = self.addcol(data, 'ho_score', 'gu_score')
        if ss1 != 0:
            return  

        self.nums +=1
        _row = []
        for i in range(len(cols)):
            _row.append(data[self._header_keys[cols[i]]])

        self.datas.append(_row)
        self.labels.append([data[len(data)-1]])

    def preparev1(self, row):
        data = []
        for i in range(len(row)):
            data.append(numpy.float32(row[i]))

        # ss1 = self.subcol(data, 'ho_score', 'gu_score') 
        # # 7091过滤
        # sb = self.getcol(data, 'size_big')
        # size = self.getcol(data, 'size')
        # if sb < 1.9 or sb > 2.0 or abs(ss1) != 1 or size < 0.1:
        #     return       

        self.nums +=1
        _row = []
        for i in range(len(cols)):
            _row.append(data[self._header_keys[cols[i]]])

        self.datas.append(_row)
        self.labels.append([data[len(data)-1]])

    def preparev2(self, row):
        self.nums +=1
        data = []
        for i in range(len(row)):
            data.append(numpy.float32(row[i]))

        wr, er, lr = self.getRatio3(data, 'first_avg_hm', 'first_avg_eq', 'first_avg_aw')
        br, sr = self.getRatio2(data, 'first_size_big', 'first_size_sma')
        wr2, er2, lr2 = self.getRatio3(data, 'avg_hm', 'avg_eq', 'avg_aw')        
        br2, sr2 = self.getRatio2(data, 'size_big', 'size_sma')
        size = self.getcol(data, 'size')
        fs = self.getcol(data, 'first_size')
        ss1 = self.subcol(data, 'ho_score', 'gu_score') # 得分差
        fss = self.subcol(data, 'ho_half_score', 'gu_half_score')
        
        _row = [wr, er, lr, br, sr, wr2, er2, lr2, br2, sr2, size, fs, ss1, fss, self.getcol(data, 'ho_score'), self.getcol(data, 'gu_score')]
        # for i in range(len(indata)):
        #     _row.append(data[self._header_keys[indata[i]]])
        
        self.datas.append(_row)

        self.labels.append([data[len(data)-1]])

    def getRatio3(self, rowdata, col1, col2, col3):
        win = rowdata[self._header_keys[col1]]
        eq = rowdata[self._header_keys[col2]]
        lost = rowdata[self._header_keys[col3]]
        if win == eq and eq == lost and lost == 0.0 :
            return 0.0, 0.0, 0.0

        #返回yxif
        D = (win*eq*lost)/(win*eq+win*lost +eq*lost)
        wr = D / win
        er = D / eq
        lr = D / lost
        return wr, er, lr

    def getRatio2(self, rowdata, col1, col2):
        win = rowdata[self._header_keys[col1]]
        lost = rowdata[self._header_keys[col2]]
        if win == lost and lost == 0.0 :
            return 0.0, 0.0

        #返回yxif
        D = (win*lost)/(win + lost)
        wr = D / win
        lr = D / lost
        return wr, lr

    def getcol(self, rowdata, col) :
        return rowdata[self._header_keys[col]]

    def addcol(self, rowdata, col1, col2):
        c1 = self._header_keys[col1]
        c2 = self._header_keys[col2]
        return rowdata[c1] + rowdata[c2]
    
    def addcol3(self, rowdata, col1, col2, col3):
        c1 = self._header_keys[col1]
        c2 = self._header_keys[col2]
        c3 = self._header_keys[col3]
        return rowdata[c1] + rowdata[c2] + rowdata[c3]

    def subcol(self, rowdata, col1, col2):
        c1 = self._header_keys[col1]
        c2 = self._header_keys[col2]
        return rowdata[c1] - rowdata[c2]

    def mulcol(self, rowdata, col1, col2):
        c1 = self._header_keys[col1]
        c2 = self._header_keys[col2]
        return rowdata[c1] * rowdata[c2]

    def mulcol3(self, rowdata, col1, col2, col3):
        c1 = self._header_keys[col1]
        c2 = self._header_keys[col2]
        c3 = self._header_keys[col3]
        return rowdata[c1] * rowdata[c2] * rowdata[c3]

    def divcol(self, rowdata, col1, col2):
        c1 = self._header_keys[col1]
        c2 = self._header_keys[col2]
        return rowdata[c1] / rowdata[c2]

    def next_batch(self, count):
        start = self._index_in_epoch
        self._index_in_epoch += count
        if self._index_in_epoch > self.nums:
            self._epochs_completed += 1
            perm = numpy.arange(self.nums)
            numpy.random.shuffle(perm)
            self.datas = numpy.array(self.datas)[perm]
            self.labels = numpy.array(self.labels)[perm]
            start = 0
            self._index_in_epoch = count
        end = self._index_in_epoch
        return self.datas[start:end], self.labels[start:end]

#定义一个隐藏层
def add_layer(inputs,in_size,out_size,activation_function = None):
    #初始化权重，一般是随机初始化
    Weights = tf.Variable(tf.random_normal([in_size,out_size]), name="W")
    #初始化偏置，一般会在0的基础上加0.1
    biases = tf.Variable(tf.zeros([1,out_size]) + 0.1, name="b")
    #计算Wx+b
    feature = tf.matmul(inputs,Weights) + biases

    #判断是否需要激活函数
    if activation_function is None:
        outputs = feature
    else:
        outputs = activation_function(feature)

    #输出激活后的隐特征
    return outputs

def train():
    dataset = DataSet("half.csv") # 数据集
    count = dataset.getDataLen() # 输入参数个数

    x = tf.placeholder(tf.float32, [None, count]) # 输入数据
    y_ = tf.placeholder(tf.float32, [None,1]) # 正确结果值

    y_hat = add_layer(x, count, 1, tf.sigmoid)
    # 两层
    # hide_count = count
    # layout1 = add_layer(x, count, hide_count, tf.tanh)   #使用tanh激活函数使模型非线性化
    # y_hat = add_layer(layout1, hide_count, 1, tf.sigmoid)#sigmoid将逻辑回归的输出概率化

    # 方差损失函数，逻辑回归不能用
    # cost = -tf.reduce_mean(tf.square(y_ - y_hat))
    # clip_by_value函数将y限制在1e-10和1.0的范围内，防止出现log0的错误，即防止梯度消失或爆发
    cross_entropy = -tf.reduce_mean(y_ * tf.log(tf.clip_by_value(y_hat, 1e-10, 1.0)) + (1-y_)*tf.log(tf.clip_by_value((1-y_hat), 1e-10, 1.0)))

    train_step = tf.train.AdamOptimizer(0.0001).minimize((cross_entropy))
    accuracy = tf.reduce_mean(tf.cast(tf.equal(tf.round(y_hat), y_), "float"))
    init = tf.global_variables_initializer()
    saver=tf.train.Saver(max_to_keep=1)

    with tf.Session() as sess:
        sess.run(init)
        for v in tf.global_variables():
	        print(v)
        min_acc = 100
        for i in range(100000):
            batch_xs,batch_ys = dataset.next_batch(500)
            sess.run(train_step, feed_dict={x: batch_xs, y_: batch_ys})
            if i % 1000 == 0:
                # 每隔一段时间计算在所有数据上的损失函数并输出
                total_cross_entropy = sess.run(cross_entropy, feed_dict={x: batch_xs, y_: batch_ys})
                print("After {} training steps(s), cross entropy on all data is {}\n".format(i, total_cross_entropy))
                if min_acc > total_cross_entropy:
                    min_acc = total_cross_entropy
                    saver.save(sess,'bet365.ckpt',global_step=i+1)


        batch_xs, batch_ys = dataset.next_batch(10000)
        match = "match {:.2f}%".format(sess.run(accuracy,feed_dict={x: batch_xs, y_: batch_ys})*100)
        print( match)

def test():
    dataset = DataSet("halftest.csv") # 数据集
    count = dataset.getDataLen() # 输入参数个数

    x = tf.placeholder(tf.float32, [None, count]) # 输入数据
    y_ = tf.placeholder(tf.float32, [None,1]) # 正确结果值

    y_hat = add_layer(x, count, 1, tf.sigmoid)
    # 两层
    # hide_count = count
    # layout1 = add_layer(x, count, hide_count, tf.tanh)   #使用tanh激活函数使模型非线性化
    # y_hat = add_layer(layout1, hide_count, 1, tf.sigmoid)#sigmoid将逻辑回归的输出概率化

    accuracy = tf.reduce_mean(tf.cast(tf.equal(tf.round(y_hat), y_), "float"))
    accuracy1 = (tf.reduce_sum(tf.multiply(tf.round(y_hat), y_)/tf.reduce_sum(tf.round(y_hat))))

    saver=tf.train.Saver()
    with tf.Session() as sess:
        model_file=tf.train.latest_checkpoint('./')
        saver.restore(sess,model_file)
        print(sess.run(tf.get_default_graph().get_tensor_by_name("W:0")))
        print(sess.run(tf.get_default_graph().get_tensor_by_name("b:0")))
        for i in range(1):
            batch_xs, batch_ys = dataset.next_batch(20)
            # result = sess.run(tf.round(y_hat), feed_dict={x: batch_xs, y_: batch_ys})
            # print(result)
            match = "match {:.2f}%".format(sess.run(accuracy,feed_dict={x: batch_xs, y_: batch_ys})*100)
            print( match)
            match = "match2 {:.2f}%\n".format(sess.run(accuracy1,feed_dict={x: batch_xs, y_: batch_ys})*100)
            print( match)


def main(argv):
    print(argv[1])
    if argv[1] == "train" :
        train()
    else:
        test()
 
if __name__ == "__main__":
    main(sys.argv)
