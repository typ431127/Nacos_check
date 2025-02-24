<template>
  <div class="service-container">
    <h1>Nacos服务实例列表</h1>
    <div class="search-bar">
      <el-input v-model="searchQuery" type="text" placeholder="输入服务名称搜索..." @input="filterServices" />
    </div>
    <div v-if="filteredServices && Object.keys(filteredServices).length > 0" class="service-list">
      <div v-for="(serviceGroup, serviceName) in filteredServices" :key="serviceName" class="service-card">
        <h2>{{ serviceName }} ({{ serviceGroup.length }}个实例)</h2>
        <div class="instance-list">
          <el-tag v-for="instance in serviceGroup" :key="instance.targets[0]" :class="{ 'containerd-true': instance.labels.containerd === 'true' }" type="primary">
            <img v-if="instance.labels.containerd === 'true'" src="@/assets/docker.svg" class="icon-docker" />{{ instance.targets[0] }}
          </el-tag>
        </div>
      </div>
    </div>
    <div v-else class="no-results">
      没有找到匹配的服务
    </div>
  </div>
</template>

<script>
export default {
  name: 'ServiceInstance',
  data() {
    return {
      services: [],
      searchQuery: '',
      filteredServices: {}
    }
  },
  computed: {
    groupedServices() {
      if (!Array.isArray(this.services)) return {}
      const groups = {}
      this.services.forEach(item => {
        const serviceName = item.labels?.service || 'unknown'
        if (!groups[serviceName]) {
          groups[serviceName] = []
        }
        groups[serviceName].push(item)
      })
      return groups
    }
  },
  async created() {
    await this.fetchData()
    console.log('分组后的数据:', this.groupedServices)
    this.filteredServices = { ...this.groupedServices }
  },
  methods: {
    async fetchData() {
      try {
        const response = await fetch('http://127.0.0.1:8099/')
        const data = await response.json()
        this.services = data
      } catch (error) {
        console.error('获取服务数据失败:', error)
        // 测试用数据
        this.services = [{"targets":["192.168.31.140:5004"],"labels":{"service":"order-server"}}, {"targets":["192.168.31.140:5009"],"labels":{"service":"push-server"}}]
      }
    },
    filterServices() {
      const query = this.searchQuery.toLowerCase().trim()
      if (!query) {
        this.filteredServices = { ...this.groupedServices }
      } else {
        this.filteredServices = Object.fromEntries(
            Object.entries(this.groupedServices).filter(([serviceName]) =>
                serviceName.toLowerCase().includes(query)
            )
        )
      }
    }
  }
}
</script>

<style scoped>
.service-container {
  max-width: 1000px;
  margin: 20px auto;
  padding: 20px;
}

.search-bar {
  margin-bottom: 20px;
}

.search-bar input {
  width: 100%;
  padding: 10px;
  font-size: 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}
.icon-docker {
  width: 14px;
  height: 14px;
  margin-right: 6px;
}
.service-list {
  /* 容器样式，无需额外调整 */
}

.service-card {
  background: #f5f5f5;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
}

h2 {
  color: #333;
  margin-bottom: 15px;
}

.instance-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.instance-item {
  background: #fff;
  padding: 8px 12px;
  border-radius: 4px;
  border: 1px solid #ddd;
  white-space: nowrap;
}

.instance-item:hover {
  background: #f0f0f0;
}

.no-results {
  text-align: center;
  color: #666;
  padding: 20px;
}
</style>