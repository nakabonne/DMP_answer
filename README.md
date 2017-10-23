# DMP_answer

## Architecture

![result](https://github.com/ryonakao/DMP_answer/blob/master/media/ArchitectureB.png)

上記アーキテクチャの"AnswerAPI"の部分。
CollectAPIは→https://github.com/ryonakao/DMP_collect

- 特定のIDFAとtimestampを指定したJSONを受け取る
- 緯度経度から最も近い駅名を返す

# Digression

当初は以下のようにジョブキューを噛ませる予定だったが、2000QPSなら使わなくても耐えることができた。

![result](https://github.com/ryonakao/DMP_answer/blob/master/media/architectureA.png)
