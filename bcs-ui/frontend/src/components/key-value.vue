<template>
    <div class="key-value">
        <div class="key-value-item">
            <span class="key desc">
                {{$t('键')}}:
                <i v-bk-tooltips="keyDesc"
                    class="ml10 bcs-icon bcs-icon-question-circle"
                    v-if="keyDesc"
                ></i>
            </span>
            <span class="value desc">
                {{$t('值')}}:
                <i v-bk-tooltips="valueDesc"
                    class="ml10 bcs-icon bcs-icon-question-circle"
                    v-if="valueDesc"
                ></i>
            </span>
        </div>
        <div v-for="(item, index) in keyValueData"
            :key="index"
            class="key-value-item"
        >
            <bcs-input :placeholder="$t('键')" :disabled="item.disabled" class="key" v-model="item.key"></bcs-input>
            <span class="equals-sign">=</span>
            <bcs-input :placeholder="item.placeholder || $t('值')" class="value" v-model="item.value"></bcs-input>
            <i class="bk-icon icon-plus-circle ml10 mr5" @click="handleAddKeyValue(index)"></i>
            <i class="bk-icon icon-minus-circle" v-if="keyValueData.length > 1" @click="handleDeleteKeyValue(index)"></i>
        </div>
        <div class="mt30">
            <bcs-button class="bcs-btn"
                theme="primary"
                :loading="loading"
                :disalbed="loading"
                @click="confirmSetLabel"
            >
                {{$t('保存')}}
            </bcs-button>
            <bcs-button class="bcs-btn" :disalbed="loading" @click="hideSetLabel">
                {{$t('取消')}}
            </bcs-button>
        </div>
    </div>
</template>
<script lang="ts">
    import { defineComponent, toRefs, watch, ref, computed } from '@vue/composition-api'

    export interface IData {
        key: string;
        value: string;
        placeholder?: any;
        disabled?: boolean;
    }
    export default defineComponent({
        props: {
            modelValue: {
                type: [Object, Array],
                default: []
            },
            valueDesc: {
                type: String,
                default: ''
            },
            keyDesc: {
                type: String,
                default: ''
            },
            loading: {
                type: Boolean,
                default: false
            }
        },
        model: {
            prop: 'modelValue',
            event: 'change'
        },
        setup (props, ctx) {
            const { modelValue } = toRefs(props)
            const keyValueData = ref<IData[]>([])
            watch(modelValue, () => {
                if (Array.isArray(modelValue.value)) {
                    keyValueData.value = modelValue.value.map((item: any) => ({
                        ...item,
                        disabled: true
                    }))
                } else {
                    keyValueData.value = Object.keys(modelValue.value).map(key => {
                        return {
                            key,
                            value: modelValue.value[key],
                            disabled: true
                        }
                    })
                }
                // 添加一组空值
                keyValueData.value.push({
                    key: '',
                    value: ''
                })
            }, { immediate: true })

            const handleAddKeyValue = (index) => {
                keyValueData.value.splice(index + 1, 0, {
                    key: '',
                    value: ''
                })
            }
            const handleDeleteKeyValue = (index) => {
                keyValueData.value.splice(index, 1)
            }
            const labels = computed(() => {
                return keyValueData.value.filter(item => !!item.key).reduce((pre, curLabelItem) => {
                    pre[curLabelItem.key] = curLabelItem.value
                    return pre
                }, {})
            })
            const checkKeyDuplicated = () => {
                const data = keyValueData.value.map(item => item.key)
                const removeDuplicateData = new Set(data)
                const result = data.filter(key => !removeDuplicateData.has(key))
                if (result.length) {
                    ctx.root.$bkMessage({
                        theme: 'error',
                        message: ctx.root.$i18n.t('键值【{key}】重复，请重新填写', { key: result[0] })
                    })
                }
                return !!result.length
            }
            const confirmSetLabel = () => {
                if (checkKeyDuplicated()) return
                ctx.emit('confirm', labels.value)
            }
            const hideSetLabel = () => {
                ctx.emit('cancel', labels.value)
            }
            return {
                keyValueData,
                confirmSetLabel,
                hideSetLabel,
                handleAddKeyValue,
                handleDeleteKeyValue
            }
        }
    })
</script>
<style lang="postcss" scoped>
.key-value-item {
    display: flex;
    align-items: center;
    height: 32px;
    line-height: 32px;
    margin-bottom: 10px;
    font-size: 14px;
    .key {
        flex: 1;
    }
    .value {
        flex: 1;
    }
    .desc {
        display: flex;
        align-items: center;
    }
    .bk-icon {
        font-size: 24px;
        color: #979bA5;
        cursor: pointer;
    }
    .equals-sign {
        color: #c3cdd7;
        margin: 0 15px;
    }
}
.bcs-btn {
    width: 86px;
}
</style>
